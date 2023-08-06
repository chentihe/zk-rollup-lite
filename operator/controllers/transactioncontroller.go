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
		HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := c.TransactionService.SendTransaction(tx); err != nil {
		HandleError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"result": "send tx sucess"})
}
