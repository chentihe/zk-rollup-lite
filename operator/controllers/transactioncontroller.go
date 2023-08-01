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
	if err := ctx.ShouldBindJSON(&tx); err != nil {
		panic(err)
	}

	savedTxs, err := c.TransactionService.SendTransaction(tx)
	if err != nil || savedTxs == math.MaxInt64 {
		panic(err)
	}

	if savedTxs == 1 {
		c.TransactionPubSub.Publish("execute roll up")
	}

	ctx.IndentedJSON(http.StatusOK, "tx finished")
}
