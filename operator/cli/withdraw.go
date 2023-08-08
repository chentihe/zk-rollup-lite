package cli

import (
	"context"
	"log"

	"github.com/chentihe/zk-rollup-lite/operator/circuits"
	"github.com/chentihe/zk-rollup-lite/operator/cmd/flags"
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/config/servicecontext"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/chentihe/zk-rollup-lite/operator/txutils"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/urfave/cli/v2"
)

func Withdraw(ctx *cli.Context, context context.Context, config *config.Config, svc *servicecontext.ServiceContext) error {
	circuitPath := config.Circuit.Path + "/withdraw"

	accountIndex := ctx.Int64(flags.AccountIndexFlag.Name)
	account := config.Accounts[accountIndex]
	signer, err := clients.NewSigner(context, account.EcdsaPrivKey, svc.EthClient)
	if err != nil {
		return err
	}

	user, err := NewUser(account)
	if err != nil {
		return err
	}

	withdrawAmount := txutils.ToWei(ctx.String(flags.AmountFlag.Name), 18)
	nullifier := babyjub.NewRandPrivKey()
	signature := user.privateKey.SignMimc7(babyjub.SkToBigInt(&nullifier))

	accountDto, err := svc.AccountService.GetAccountByPublicKey(user.PublicKey.String())
	if err != nil {
		return err
	}

	mtProof, err := svc.AccountTree.GenerateCircomVerifierProof(accountDto)
	if err != nil {
		return err
	}

	withdrawInputs := &circuits.WithdrawInputs{
		Account:        accountDto,
		Nullifier:      babyjub.SkToBigInt(&nullifier),
		Signature:      signature,
		WithdrawAmount: withdrawAmount,
		MTProof:        mtProof,
	}

	circuitInput, err := withdrawInputs.InputsMarshal()
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

	var withdrawOutputs circuits.WithdrawOutputs
	if err = withdrawOutputs.OutputsUnmarshal(proof); err != nil {
		return err
	}

	data, err := svc.Abi.Pack("withdraw", withdrawAmount, withdrawOutputs.Proof.A, withdrawOutputs.Proof.B, withdrawOutputs.Proof.C, withdrawOutputs.PublicSignals)
	if err != nil {
		log.Printf("Cannot pack withdraw call data: %v\n", err)
	}

	tx, err := signer.GenerateLegacyTx(svc.RollUpAddress, data, nil)
	if err != nil {
		return err
	}

	signTx, err := signer.SignTx(tx)
	if err != nil {
		return err
	}

	if err = svc.EthClient.SendTransaction(context, signTx); err != nil {
		return err
	}
	log.Printf("Withdraw success: %s\n", signTx.Hash().Hex())

	return nil
}
