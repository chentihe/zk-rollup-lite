package controllers

import (
	"net/http"

	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/chentihe/zk-rollup-lite/operator/txutils"
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
	var tx *txutils.TransactionInfo
	if err := ctx.ShouldBindJSON(&tx); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, err)
	}

	if err := c.TransactionService.SendTransaction(tx); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}

	ctx.IndentedJSON(http.StatusCreated, "send tx sucess")
}
