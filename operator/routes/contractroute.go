package routes

import (
	"github.com/chentihe/zk-rollup-lite/operator/controllers"
	"github.com/gin-gonic/gin"
)

func AddContractRoutes(v1 *gin.RouterGroup, contractController *controllers.ContractController) {
	contractGroup := v1.Group("/contract")
	contractGroup.GET("/users/:id", contractController.GetAccountByIndex)
	contractGroup.GET("/balance", contractController.GetContractBalance)
	contractGroup.GET("/root", contractController.GetStateRoot)
}
