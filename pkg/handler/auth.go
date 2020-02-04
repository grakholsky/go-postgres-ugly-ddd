package handler

import (
	"go-postgres/pkg/model"
	"go-postgres/pkg/route"
	"go-postgres/pkg/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type Auth struct {
	Router *route.Router
}

func (h *Auth) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			return
		}

		token, err := h.parseHeader(header)
		if err != nil {
			c.AbortWithStatusJSON(400, handleError(err))
			return
		}

		claims, err := h.Router.Jwt.ParseClaims(token)
		if err != nil {
			c.AbortWithStatusJSON(403, handleError(err))
			return
		}

		user, err := h.Router.Account.GetByID(claims.UserID)
		if err != nil {
			if service.IsNotFoundError(err) {
				c.AbortWithStatusJSON(404, handleError(fmt.Errorf("account not found")))
				return
			}
			c.AbortWithStatusJSON(500, handleError(err))
			return
		}

		err = h.Router.Jwt.ParseToken(token, user.Salt)
		if err != nil {
			c.AbortWithStatusJSON(403, handleError(err))
			return
		}

		c.Set("User", user)
	}
}

func (h *Auth) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		roles := []*model.Role{{Name: "anonymous"}}

		user, ok := c.Value("User").(*model.User)
		if ok {
			userRoles, err := h.Router.Role.GetListByUserID(user.ID)
			if err != nil {
				c.AbortWithStatusJSON(500, handleError(err))
				return
			}
			roles = append(roles, userRoles...)
			c.Set("UserRoles", userRoles)
		}

		for _, role := range roles {
			if h.Router.Policy.Check(role.Name, c.Request.URL.Path, c.Request.Method) {
				return
			}
		}
		c.AbortWithStatusJSON(403, handleError(fmt.Errorf("access denied")))
	}
}

func (h *Auth) parseHeader(header string) (string, error) {
	splits := strings.SplitN(header, " ", 2)
	if len(splits) < 2 {
		return "", fmt.Errorf("incorrect authorization header")
	}
	if strings.ToLower(splits[0]) != strings.ToLower("bearer") {
		return "", fmt.Errorf("incorrect bearer strategy")
	}
	return splits[1], nil
}
