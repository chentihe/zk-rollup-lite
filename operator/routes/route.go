package routes

import (
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(router *gin.Engine, svc *config.ServiceContext) {
	v1 := router.Group("/api/v1")
	AddAccountRoutes(v1, svc.AccountController)
	AddTransactionRoutes(v1, svc.TransactionController)
}
