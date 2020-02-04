package main

import (
	"github.com/gin-gonic/gin"
	"go-postgres/pkg/handler"
	"go-postgres/pkg/repository"
	"go-postgres/pkg/route"
	"go-postgres/pkg/service"
	"os"
)

var config *route.RouterConfig

func init() {
	config = &route.RouterConfig{
		PolicyConfig: &service.PolicyConfig{
			CasbinModelPath:  os.Getenv("CASBIN_MODEL_PATH"),
			CasbinPolicyPath: os.Getenv("CASBIN_POLICY_PATH"),
		},
		PostgresConfig: &repository.PostgresConfig{
			URI: os.Getenv("POSTGRES_URI"),
		},
	}
}

func main() {
	postgresDb := repository.ConnectPostgres(config.PostgresConfig)
	repoRole := repository.NewRole(postgresDb)
	svcAccount := service.NewAccount(repository.NewUser(postgresDb), repoRole)

	router := &route.Router{
		Config:  config,
		Engine:  gin.New(),
		Account: svcAccount,
		Jwt:     service.NewJwt(svcAccount),
		Policy:  service.NewPolicy(config.PolicyConfig),
		Role:    service.NewRole(repoRole),
	}

	handlerAccount := &handler.Account{Router: router}
	handlerAuth := &handler.Auth{Router: router}
	router.Engine.Use(gin.Logger(), gin.Recovery(), handlerAuth.Authenticate(), handlerAuth.Authorize())
	base := router.Engine.Group("/api/v1")
	base.Group("/account").GET("/me", handlerAccount.Me)
	base.Group("/account").POST("/sign-in", handlerAccount.SignIn)
	base.Group("/account").POST("/sign-up", handlerAccount.SignUp)
	router.Engine.Run() // listen and serve on 0.0.0.0:8080
}
