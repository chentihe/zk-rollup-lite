package daos

import "fmt"

var (
	ErrSqlOperation    = fmt.Errorf("Sql operation error")
	ErrAccountNotFound = fmt.Errorf("Account not found")
)
