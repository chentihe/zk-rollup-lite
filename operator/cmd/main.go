package main

import (
	"fmt"

	accounttree "github.com/chentihe/zk-rollup-lite/operator/tree"
)

func main() {
	mt, err := accounttree.InitAccountTree()
	if err != nil {
		panic(fmt.Sprintf("cannot create merkletree, %v\n", err))
	}

	fmt.Printf("merkle tree init success, #%v", mt)
}
