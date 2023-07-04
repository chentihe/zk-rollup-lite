package main

import (
	"fmt"

	"github.com/chentihe/zk-rollup-lite/operator/accounttree"
)

func main() {
	mt, err := accounttree.InitMerkleTree()
	if err != nil {
		panic(fmt.Sprintf("cannot create merkletree, %v\n", err))
	}

	fmt.Printf("merkle tree init success, #%v", mt)
}
