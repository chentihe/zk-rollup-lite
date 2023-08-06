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

// TODO: write a handle err func
func (c *AccountController) GetAccountByIndex(ctx *gin.Context) {
	id := ctx.Param("id")

	index, err := strconv.Atoi(id)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := c.AccountService.GetAccountByIndex(int64(index))
	if err != nil {
		HandleError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, res)
}
