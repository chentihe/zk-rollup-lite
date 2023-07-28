package cli

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/chentihe/zk-rollup-lite/operator/circuits"
	"github.com/chentihe/zk-rollup-lite/operator/cmd/flags"
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/config/servicecontext"
	"github.com/chentihe/zk-rollup-lite/operator/daos"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/urfave/cli/v2"
)

func Deposit(ctx *cli.Context, context context.Context, config *config.Config, svc *servicecontext.ServiceContext) error {
	circuitPath := config.Circuit.Path + "/deposit"

	accountIndex := ctx.Int64(flags.AccountIndexFlag.Name)
	account := config.Accounts[accountIndex]
	signer, err := clients.NewSigner(context, account.EcdsaPrivKey, svc.EthClient)
	if err != nil {
		return err
	}

	user, err := NewUser(account.EddsaPrivKey)
	if err != nil {
		return err
	}

	depositAmount := ToWei(ctx.String(flags.DepositAmountFlag.Name), 18)

	depositInputs := &circuits.DepositInputs{
		DepositAmount: depositAmount,
	}

	var mtProof *merkletree.CircomProcessorProof

	comp := user.PublicKey.String()

	accountDto, err := svc.AccountService.GetAccountByPublickKey(comp)
	if err == daos.ErrAccountNotFound {
		userIndex, err := svc.AccountService.GetCurrentAccountIndex()
		if err != nil {
			return err
		}

		accountDto = &models.AccountDto{
			AccountIndex: userIndex,
			PublicKey:    comp,
			Balance:      depositAmount,
			Nonce:        0,
		}

		mtProof, err = svc.AccountTree.AddAccount(accountDto)
		if err != nil {
			return err
		}
	} else {
		// zkp new root should be the new state root cannot use mock merkle tree proof
		// mock update to get the circuit processor proof
		accountDto.Balance = new(big.Int).Add(accountDto.Balance, depositAmount)
		mtProof, err = svc.AccountTree.UpdateAccount(accountDto)
		if err != nil {
			return err
		}
	}

	depositInputs.Account = accountDto
	depositInputs.MTProof = mtProof

	circuitInput, err := depositInputs.InputsMarshal()
	if err != nil {
		return err
	}

	proof, err := circuits.GenerateGroth16Proof(circuitInput, circuitPath)
	if err != nil {
		return err
	}

	if err = circuits.VerifierGroth16(proof, circuitPath); err != nil {
		return err
	}

	var depositOutputs circuits.DepositOutputs
	if err = depositOutputs.OutputsUnmarshal(proof); err != nil {
		return err
	}

	data, err := svc.Abi.Pack("deposit", depositOutputs.Proof.A, depositOutputs.Proof.B, depositOutputs.Proof.C, depositOutputs.PublicSignals)
	if err != nil {
		fmt.Printf("Cannot pack rollup call data: %v", err)
	}

	rollupAddress := common.HexToAddress(config.SmartContract.Address)
	tx, err := signer.GenerateDynamicTx(&rollupAddress, data, depositAmount)
	if err != nil {
		return err
	}

	signTx, err := signer.SignTx(tx)
	if err != nil {
		return err
	}

	rawTxBytes, err := signTx.MarshalBinary()
	if err != nil {
		return err
	}

	depositInfo := txmanager.DepositInfo{
		AccountIndex:  accountIndex,
		PublicKey:     user.PublicKey.String(),
		DepositAmount: depositAmount,
		SignedTxHash:  hex.EncodeToString(rawTxBytes),
		ZkProof:       proof,
	}

	requestBody, err := json.Marshal(depositInfo)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8000/api/v1/deposit", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
