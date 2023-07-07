package accounttree

import (
	"context"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
)

const mtDepth = 6

type AccountTree struct {
	MT *merkletree.MerkleTree
}

// TODO: check which one is better? memory or postgresdb
func InitAccountTree() (*AccountTree, error) {
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

	return &AccountTree{mt}, nil
}

func (accountTree *AccountTree) UpdateAccountTree(account *models.AccountModel) error {
	context := context.Background()

	leaf, err := GenerateAccountLeaf(account)
	if err != nil {
		return err
	}

	if _, err := accountTree.MT.Update(context, big.NewInt(account.AccountIndex), leaf); err != nil {
		return err
	}

	return nil
}
