package controllers

import (
	"net/http"

	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/gin-gonic/gin"
)

type AccountController struct {
	AccountService *services.AccountService
}

func NewAccountController(accountService *services.AccountService) *AccountController {
	return &AccountController{
		AccountService: accountService,
	}
}

func (c *AccountController) GetAccountByIndex(ctx *gin.Context) {
	var index int64
	if err := ctx.ShouldBindQuery(&index); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, err)
	}

	res, err := c.AccountService.GetAccountByIndex(index)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, err)
	}

	ctx.IndentedJSON(http.StatusOK, res)
}
