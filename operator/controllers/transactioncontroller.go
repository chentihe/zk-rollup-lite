package controllers

import (
	"math"
	"net/http"

	"github.com/chentihe/zk-rollup-lite/operator/pubsubs"
	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	TransactionService *services.TransactionService
	TransactionPubSub  pubsubs.Subscriber
}

func NewTransactionController(transactionService *services.TransactionService, pubsub *pubsubs.Subscriber) *TransactionController {
	return &TransactionController{
		TransactionService: transactionService,
		TransactionPubSub:  *pubsub,
	}
}

func (c *TransactionController) SendTransaction(ctx *gin.Context) {
	var tx *txmanager.TransactionInfo
	if err := ctx.ShouldBindQuery(&tx); err != nil {
		panic(err)
	}

	savedTxs, err := c.TransactionService.SendTransaction(tx)
	if err != nil || savedTxs == math.MaxInt64 {
		panic(err)
	}

	if savedTxs == -1 {
		c.TransactionPubSub.Publish("execute roll up")
	}

	ctx.IndentedJSON(http.StatusOK, "tx finished")
}

func (c *TransactionController) Deposit(ctx *gin.Context) {
	var deposit *txmanager.DepositInfo
	if err := ctx.ShouldBindQuery(&deposit); err != nil {
		panic(err)
	}

	err := c.TransactionService.Deposit(deposit)
	if err != nil {
		panic(err)
	}

	c.TransactionPubSub.Publish(deposit.SignedTxHash)

	ctx.IndentedJSON(http.StatusOK, "deposit finished")
}

func (c *TransactionController) Withdraw(ctx *gin.Context) {
	var withdraw *txmanager.WithdrawInfo
	if err := ctx.ShouldBindQuery(&withdraw); err != nil {
		panic(err)
	}

	err := c.TransactionService.Withdraw(withdraw)
	if err != nil {
		panic(err)
	}

	c.TransactionPubSub.Publish(withdraw.SignedTxHash)

	ctx.IndentedJSON(http.StatusOK, "withdraw finished")
}
