package controllers

import (
	"net/http"

	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	TransactionService *services.TransactionService
}

func NewTransactionController(transactionService *services.TransactionService) *TransactionController {
	return &TransactionController{
		TransactionService: transactionService,
	}
}

func (c *TransactionController) SendTransaction(ctx *gin.Context) {
	var tx txmanager.TransactionInfo
	if err := ctx.ShouldBindQuery(&tx); err != nil {
		panic(err)
	}

	if err := c.TransactionService.SendTransaction(&tx); err != nil {
		panic(err)
	}

	ctx.IndentedJSON(http.StatusOK, "tx finished")
}

func (c *TransactionController) Deposit(ctx *gin.Context) {
	var deposit txmanager.DepositInfo
	if err := ctx.ShouldBindQuery(&deposit); err != nil {
		panic(err)
	}

	if err := c.TransactionService.Deposit(deposit); err != nil {
		panic(err)
	}

	ctx.IndentedJSON(http.StatusOK, "deposit finished")
}

func (c *TransactionController) Withdraw(ctx *gin.Context) {
	var withdraw txmanager.WithdrawInfo
	if err := ctx.ShouldBindQuery(&withdraw); err != nil {
		panic(err)
	}

	if err := c.TransactionService.Withdraw(withdraw); err != nil {
		panic(err)
	}

	ctx.IndentedJSON(http.StatusOK, "withdraw finished")
}
