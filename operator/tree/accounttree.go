package tree

import (
	"context"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	sql "github.com/iden3/go-merkletree-sql/db/pgx/v5"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

const mtDepth = 5

type AccountTree struct {
	MT      *merkletree.MerkleTree
	PgxPool *pgxpool.Pool
}

func InitAccountTree(context context.Context, ethClient *ethclient.Client, abi *abi.ABI, contractAddress *common.Address, config *config.Postgres) (*AccountTree, error) {
	pool, err := pgxpool.New(context, config.Url())
	if err != nil {
		return nil, err
	}

	treeStorage := sql.NewSqlStorage(pool, 1)

	mt, err := merkletree.NewMerkleTree(context, treeStorage, mtDepth)
	if err != nil {
		return nil, err
	}

	return &AccountTree{
		MT:      mt,
		PgxPool: pool,
	}, nil
}

func (accountTree *AccountTree) RestoreTree() (*merkletree.MerkleTree, error) {
	treeStroage := sql.NewSqlStorage(accountTree.PgxPool, 1)
	return merkletree.NewMerkleTree(context.Background(), treeStroage, mtDepth)
}

func (accountTree *AccountTree) UpdateAccountTree(accountDto *models.AccountDto) (*merkletree.CircomProcessorProof, error) {
	context := context.Background()

	mt, err := accountTree.RestoreTree()
	if err != nil {
		return nil, err
	}

	leaf, err := GenerateAccountLeaf(accountDto)
	if err != nil {
		return nil, err
	}

	proof, err := mt.Update(context, big.NewInt(accountDto.AccountIndex), leaf)
	if err != nil {
		return nil, err
	}

	return proof, nil
}

func (accountTree *AccountTree) GetPathByAccount(account *models.AccountDto) ([]*merkletree.Hash, error) {
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

	mt, err := accountTree.RestoreTree()
	if err != nil {
		return nil, err
	}

	return mt.AddAndGetCircomProof(context, big.NewInt(key), value)
}

func (accountTree *AccountTree) Delete(key int64) error {
	context := context.Background()

	mt, err := accountTree.RestoreTree()
	if err != nil {
		return err
	}

	hashKey, err := merkletree.NewHashFromBigInt(big.NewInt(key))
	if err != nil {
		return err
	}

	if _, err = mt.DumpLeafs(context, hashKey); err != nil {
		return err
	}

	return nil
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
