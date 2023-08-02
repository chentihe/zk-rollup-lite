package routes

import (
	"github.com/chentihe/zk-rollup-lite/operator/controllers"
	"github.com/gin-gonic/gin"
)

func AddAccountRoutes(v1 *gin.RouterGroup, accountController *controllers.AccountController) {
	accountGroup := v1.Group("/accounts")
	accountGroup.GET(":id", accountController.GetAccountByIndex)
}
