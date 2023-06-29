package database

import (
	"context"
	"errors"

	sql "github.com/iden3/go-merkletree-sql/db/pgx/v5"
	"github.com/iden3/go-merkletree-sql/v2"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitMerkleTree() (*merkletree.MerkleTree, error) {
	// TODO: move to env
	urlExample := "postgres://username:password@localhost:5432/database-name"
	mtDepth := 6
	mtId := uint64(1)

	// TODO: pgxPool & context should move to main.go
	// pass into this func as an arg
	ctx := context.Background()
	pgxPool, err := pgxpool.New(ctx, urlExample)
	if err != nil {
		return nil, errors.New("unable to connect to the database")
	}
	defer pgxPool.Close()

	treeStorage := sql.NewSqlStorage(pgxPool, mtId)

	mt, err := merkletree.NewMerkleTree(ctx, treeStorage, mtDepth)
	if err != nil {
		return nil, err
	}

	return mt, nil
}
