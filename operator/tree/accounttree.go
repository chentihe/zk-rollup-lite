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
	context context.Context
	pgxPool *pgxpool.Pool
}

// sparse merkle tree reserves the index 0
// don't insert leaf into index 0
func InitAccountTree(context context.Context, ethClient *ethclient.Client, abi *abi.ABI, contractAddress *common.Address, config *config.Postgres) (*AccountTree, error) {
	pool, err := pgxpool.New(context, config.Url())
	if err != nil {
		return nil, err
	}

	return &AccountTree{
		context: context,
		pgxPool: pool,
	}, nil
}

// RestoreTree get the lastest status of the mt
// to make sure the consistency of the mt
func (accountTree *AccountTree) RestoreTree() (*merkletree.MerkleTree, error) {
	treeStroage := sql.NewSqlStorage(accountTree.pgxPool, 1)
	return merkletree.NewMerkleTree(accountTree.context, treeStroage, mtDepth)
}

func (accountTree *AccountTree) AddAccount(accountDto *models.AccountDto) (proof *merkletree.CircomProcessorProof, err error) {
	mt, err := accountTree.RestoreTree()
	if err != nil {
		return nil, err
	}

	key := big.NewInt(accountDto.AccountIndex)

	leaf, err := GenerateAccountLeaf(accountDto)
	if err != nil {
		return nil, err
	}

	return mt.AddAndGetCircomProof(accountTree.context, key, leaf)
}

func (accountTree *AccountTree) UpdateAccount(accountDto *models.AccountDto) (*merkletree.CircomProcessorProof, error) {
	mt, err := accountTree.RestoreTree()
	if err != nil {
		return nil, err
	}

	leaf, err := GenerateAccountLeaf(accountDto)
	if err != nil {
		return nil, err
	}

	// circuit processor proof
	proof, err := mt.Update(accountTree.context, big.NewInt(accountDto.AccountIndex), leaf)
	if err != nil {
		return nil, err
	}

	return proof, nil
}

func (accountTree *AccountTree) GetPathByAccount(account *models.AccountDto) ([]*merkletree.Hash, error) {
	mt, err := accountTree.RestoreTree()
	if err != nil {
		return nil, err
	}

	index := account.AccountIndex
	_, _, siblings, err := mt.Get(accountTree.context, big.NewInt(index))
	if err != nil {
		return nil, err
	}

	// fill the empty path
	siblings = merkletree.CircomSiblingsFromSiblings(siblings, mtDepth)

	return siblings, nil
}

// GenerateCircomVerifierProof to proof that the account is on the merkle tree
func (accountTree *AccountTree) GenerateCircomVerifierProof(account *models.AccountDto) (*merkletree.CircomVerifierProof, error) {
	mt, err := accountTree.RestoreTree()
	if err != nil {
		return nil, err
	}

	index := account.AccountIndex
	proof, err := mt.GenerateCircomVerifierProof(accountTree.context, big.NewInt(index), nil)
	if err != nil {
		return nil, err
	}

	return proof, nil
}

func (accountTree *AccountTree) GetRoot() *merkletree.Hash {
	mt, _ := accountTree.RestoreTree()
	return mt.Root()
}
