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
	"github.com/chentihe/zk-rollup-lite/operator/tree"
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
		Root:          svc.AccountTree.GetRoot(),
		DepositAmount: depositAmount,
	}

	var mtProof *merkletree.CircomProcessorProof

	accountDto, err := svc.AccountService.GetAccountByIndex(accountIndex)

	// TODO: will occur err if the account exists
	if err == daos.ErrAccountNotFound {
		userIndex, err := svc.AccountService.GetCurrentAccountIndex()
		if err != nil {
			return err
		}

		accountDto = &models.AccountDto{
			AccountIndex: userIndex,
			PublicKey:    user.PublicKey.String(),
			Balance:      big.NewInt(0),
			Nonce:        0,
		}

		leaf, err := tree.GenerateAccountLeaf(accountDto)
		if err != nil {
			return err
		}

		mtProof, err = svc.AccountTree.AddAndGetCircomProof(userIndex, leaf)
		if err != nil {
			return err
		}
	} else {
		// mock update to get the circuit processor proof
		mtProof, err = svc.AccountTree.UpdateAccountTree(accountDto)
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
