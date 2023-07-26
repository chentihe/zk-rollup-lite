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
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/urfave/cli/v2"
)

func Withdraw(ctx *cli.Context, context context.Context, config *config.Config, svc *servicecontext.ServiceContext) error {
	signer, err := clients.NewSigner(context, config.Sender.PrivateKey, svc.EthClient)
	if err != nil {
		return err
	}

	// TODO: need to set l2 priv key into yaml
	l2PrivateKey := babyjub.NewRandPrivKey()
	l2PublicKey := l2PrivateKey.Public()

	accountIndex := ctx.Int64(flags.AccountIndexFlag.Name)
	withdrawAmount := new(big.Int)
	withdrawAmount, ok := withdrawAmount.SetString(ctx.String(flags.WithdrawAmountFlag.Name), 10)
	if !ok {
		return fmt.Errorf("cannot convert deposit amount to big int")
	}

	signature := l2PrivateKey.SignMimc7(withdrawAmount)

	withdrawInputs := &circuits.WithdrawInputs{
		Root:           svc.AccountTree.GetRoot(),
		WithdrawAmount: withdrawAmount,
		Signature:      signature,
	}

	var mtProof *merkletree.CircomProcessorProof

	accountDto, err := svc.AccountService.GetAccountByIndex(accountIndex)

	// TODO: will occur err if the account exists
	// merkle tree issue, if don't save merkle tree in db,
	// need to figure out how to recover the merkle tree
	if err == daos.ErrAccountNotFound {
		userIndex, err := svc.AccountService.GetCurrentAccountIndex()
		if err != nil {
			return err
		}

		accountDto = &models.AccountDto{
			AccountIndex: userIndex,
			PublicKey:    l2PublicKey.String(),
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

	withdrawInputs.Account = accountDto
	withdrawInputs.MTProof = mtProof

	circuitInput, err := withdrawInputs.InputsMarshal()
	if err != nil {
		return err
	}

	proof, err := circuits.GenerateGroth16Proof(circuitInput, circuitPath+"/withdraw")
	if err != nil {
		return err
	}

	if err = circuits.VerifierGroth16(proof, circuitPath+"/withdraw"); err != nil {
		return err
	}

	var withdrawOutputs circuits.WithdrawOutputs
	if err = withdrawOutputs.OutputUnmarshal(proof); err != nil {
		return err
	}

	data, err := svc.Abi.Pack("withdraw", withdrawAmount, withdrawOutputs.Proof.A, withdrawOutputs.Proof.B, withdrawOutputs.Proof.C, withdrawOutputs.PublicSignals)
	if err != nil {
		fmt.Printf("Cannot pack rollup call data: %v", err)
	}

	rollupAddress := common.HexToAddress(config.SmartContract.Address)
	tx, err := signer.GenerateDynamicTx(&rollupAddress, data, withdrawAmount)
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

	withdrawInfo := txmanager.WithdrawInfo{
		AccountIndex:   accountIndex,
		PublicKey:      l2PublicKey.String(),
		Signature:      signature,
		WithdrawAmount: withdrawAmount,
		SignedTxHash:   hex.EncodeToString(rawTxBytes),
	}

	requestBody, err := json.Marshal(withdrawInfo)
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
