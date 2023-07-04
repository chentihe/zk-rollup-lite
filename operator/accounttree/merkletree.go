package accounttree

import (
	"context"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
)

const mtDepth = 6

// TODO: check which one is better? memory or postgresdb
func InitMerkleTree() (*merkletree.MerkleTree, error) {
	// TODO: move to env
	// urlExample := "postgres://username:password@localhost:5432/database-name"
	// mtId := uint64(1)

	// TODO: pgxPool & context should move to main.go
	// pass into this func as an arg
	ctx := context.Background()
	// pgxPool, err := pgxpool.New(ctx, urlExample)
	// if err != nil {
	// 	return nil, errors.New("unable to connect to the database")
	// }
	// defer pgxPool.Close()

	// treeStorage := sql.NewSqlStorage(pgxPool, mtId)
	treeStorage := memory.NewMemoryStorage()

	mt, err := merkletree.NewMerkleTree(ctx, treeStorage, mtDepth)
	if err != nil {
		return nil, err
	}

	return mt, nil
}

func GenerateAccountLeaf(account *models.AccountModel) (*big.Int, error) {
	publicKey, err := txmanager.DecodePublicKeyFromString(account.PublicKey)
	if err != nil {
		return nil, err
	}

	hashedLeaf, err := poseidon.Hash([]*big.Int{
		publicKey.X,
		publicKey.Y,
		account.Balance,
		big.NewInt(account.Nonce),
	})
	if err != nil {
		return nil, err
	}

	return hashedLeaf, nil
}
