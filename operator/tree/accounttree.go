package tree

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

func (accountTree *AccountTree) UpdateAccountTree(account *models.Account) (*merkletree.CircomProcessorProof, error) {
	context := context.Background()

	leaf, err := GenerateAccountLeaf(account)
	if err != nil {
		return nil, err
	}

	proof, err := accountTree.MT.Update(context, big.NewInt(account.AccountIndex), leaf)
	if err != nil {
		return nil, err
	}

	return proof, nil
}

func (accountTree *AccountTree) GetPathByAccount(account *models.Account) ([]*merkletree.Hash, error) {
	context := context.Background()

	index := account.AccountIndex

	_, _, siblings, err := accountTree.MT.Get(context, big.NewInt(index))
	if err != nil {
		return nil, err
	}

	// fill the empty path
	siblings = merkletree.CircomSiblingsFromSiblings(siblings, mtDepth)

	return siblings, nil
}

func (accountTree *AccountTree) GetRoot() *merkletree.Hash {
	return accountTree.MT.Root()
}

func (accountTree *AccountTree) Add(key int64, value *big.Int) error {
	context := context.Background()
	return accountTree.MT.Add(context, big.NewInt(key), value)
}

func (accountTree *AccountTree) AddAndGetCircomProof(key int64, value *big.Int) (proof *merkletree.CircomProcessorProof, err error) {
	context := context.Background()
	return accountTree.MT.AddAndGetCircomProof(context, big.NewInt(key), value)
}

func (accountTree *AccountTree) GenerateProof(key *big.Int) (proof *merkletree.CircomVerifierProof, err error) {
	context := context.Background()
	root := accountTree.GetRoot()
	proof, err = accountTree.MT.GenerateCircomVerifierProof(context, key, root)
	if err != nil {
		return nil, err
	}

	return proof, nil
}
