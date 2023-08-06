package controllers

import (
	"math/big"
	"net/http"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/gin-gonic/gin"
)

type ContractController struct {
	ContractService *services.ContractService
}

func NewContractController(contractService *services.ContractService) *ContractController {
	return &ContractController{
		ContractService: contractService,
	}
}

func (c *ContractController) GetUserByIndex(ctx *gin.Context) {
	id := ctx.Param("id")

	index, err := strconv.Atoi(id)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := c.ContractService.GetUserByIndex(big.NewInt(int64(index)))
	if err != nil {
		HandleError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, res)
}

func (c *ContractController) GetContractBalance(ctx *gin.Context) {
	res, err := c.ContractService.GetContractBalance()
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, res)
}

func (c *ContractController) GetStateRoot(ctx *gin.Context) {
	res, err := c.ContractService.GetStateRoot()
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, res)
}
