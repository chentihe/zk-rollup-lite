package routes

import (
	"github.com/chentihe/zk-rollup-lite/operator/controllers"
	"github.com/gin-gonic/gin"
)

func AddTransactionRoutes(v1 *gin.RouterGroup, transactionController *controllers.TransactionController) {
	v1.POST("/send", transactionController.SendTransaction)
}
