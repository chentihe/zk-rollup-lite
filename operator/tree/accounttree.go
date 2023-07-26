package tree

import (
	"context"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/layer1"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
)

const mtDepth = 5

type AccountTree struct {
	MT *merkletree.MerkleTree
}

var depositHash = crypto.Keccak256Hash([]byte("Deposit((uint256,uint256,uint256,uint256,uint256))"))

func InitAccountTree(context context.Context, ethClient *ethclient.Client, abi *abi.ABI, contractAddress *common.Address) (*AccountTree, error) {
	treeStorage := memory.NewMemoryStorage()
	mt, err := merkletree.NewMerkleTree(context, treeStorage, mtDepth)
	if err != nil {
		return nil, err
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			*contractAddress,
		},
		Topics: [][]common.Hash{{
			depositHash,
		}},
	}

	logs, err := ethClient.FilterLogs(context, query)
	if err != nil {
		return nil, err
	}

	for _, vLog := range logs {
		var deposit layer1.Deposit
		if err := abi.UnpackIntoInterface(&deposit, "Deposit", vLog.Data); err != nil {
			return nil, err
		}

		user := deposit.User
		publicKey := babyjub.PublicKey(babyjub.Point{X: user.PublicKeyX, Y: user.PublicKeyY})

		accountDto := &models.AccountDto{
			AccountIndex: user.Index.Int64(),
			PublicKey:    publicKey.String(),
			Nonce:        user.Nonce.Int64(),
		}

		accountLeaf, err := GenerateAccountLeaf(accountDto)
		if err != nil {
			return nil, err
		}

		mt.Add(context, user.Index, accountLeaf)
	}

	return &AccountTree{mt}, nil
}

func (accountTree *AccountTree) UpdateAccountTree(accountDto *models.AccountDto) (*merkletree.CircomProcessorProof, error) {
	context := context.Background()

	leaf, err := GenerateAccountLeaf(accountDto)
	if err != nil {
		return nil, err
	}

	proof, err := accountTree.MT.Update(context, big.NewInt(accountDto.AccountIndex), leaf)
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
