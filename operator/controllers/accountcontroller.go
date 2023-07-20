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
		panic(err)
	}

	res, err := c.AccountService.GetAccountByIndex(index)
	if err != nil {
		panic(err)
	}

	ctx.IndentedJSON(http.StatusOK, res)
}

func (c *AccountController) GetCurrentAccountIndex(ctx *gin.Context) {
	res, err := c.AccountService.GetCurrentAccountIndex()
	if err != nil {
		panic(err)
	}

	ctx.IndentedJSON(http.StatusOK, res)
}
