package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	dir  = "../files"
	pass = os.Getenv("KEY_PASSWORD")
)

func main() {
	abs, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	keystorePath := abs + "/keystore"

	files, err := os.ReadDir(keystorePath)
	if err != nil {
		panic(err)
	}

	for i, file := range files {
		fmt.Printf("Start Output PrivateKey #%d: %s\n", i+1, file.Name())
		filePath := keystorePath + "/" + file.Name()
		outPath := fmt.Sprintf("%s/keys/%s.hex", abs, file.Name())
		baseDir := path.Dir(outPath)
		os.MkdirAll(baseDir, 0755)

		keyJson, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		key, err := keystore.DecryptKey(keyJson, pass)
		if err != nil {
			panic(err)
		}

		err = crypto.SaveECDSA(outPath, key.PrivateKey)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Finish Output PrivateKey #%d: %s\n", i+1, file.Name())
	}
}
