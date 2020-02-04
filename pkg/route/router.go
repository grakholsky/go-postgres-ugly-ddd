package route

import (
	"github.com/gin-gonic/gin"
	"go-postgres/pkg/repository"
	"go-postgres/pkg/service"
)

type RouterConfig struct {
	PolicyConfig   *service.PolicyConfig
	PostgresConfig *repository.PostgresConfig
}

type Router struct {
	Config  *RouterConfig
	Engine  *gin.Engine
	Policy  *service.Policy
	Role    *service.Role
	Account *service.Account
	Jwt     *service.Jwt
}
