package controllers

import (
	"net/http"
	"strconv"

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
	id := ctx.Param("id")

	index, err := strconv.Atoi(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, err)
	}

	res, err := c.AccountService.GetAccountByIndex(int64(index))
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, err)
	}

	ctx.IndentedJSON(http.StatusOK, res)
}
