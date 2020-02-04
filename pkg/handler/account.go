package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-postgres/pkg/dto"
	"go-postgres/pkg/model"
	"go-postgres/pkg/route"
	"go-postgres/pkg/service"
)

type Account struct {
	Router *route.Router
}

func (h *Account) Me(c *gin.Context) {
	c.JSON(200, gin.H{"data": c.Value("User").(*model.User)})
}

func (h *Account) SignIn(c *gin.Context) {
	var dtoBody dto.SignInUser
	err := c.ShouldBindJSON(&dtoBody)
	if err != nil {
		c.JSON(400, handleError(dto.Errors(dtoBody, err, dto.BindJSON)...))
		return
	}

	user, err := h.Router.Account.GetByEmail(dtoBody.Email)
	if err != nil {
		if service.IsNotFoundError(err) {
			c.JSON(400, handleError(fmt.Errorf("your email or password do not match")))
			return
		}
		c.JSON(500, handleError(err))
		return
	}

	valid, err := h.Router.Account.CheckPassword(user.ID, dtoBody.Password)
	if err != nil {
		c.JSON(500, handleError(err))
		return
	}
	if !valid {
		c.JSON(400, handleError(fmt.Errorf("your email or password do not match")))
		return
	}

	token, err := h.Router.Jwt.GenToken(user.ID, user.Salt)
	if err != nil {
		c.JSON(500, handleError(err))
		return
	}
	c.JSON(200, gin.H{"token": token})
}

func (h *Account) SignUp(c *gin.Context) {
	var dtoBody dto.SignUpUser
	err := c.ShouldBind(&dtoBody)
	if err != nil {
		c.JSON(400, handleError(dto.Errors(dtoBody, err, dto.BindFORM)...))
		return
	}

	_, err = h.Router.Account.GetByEmail(dtoBody.Email)
	if err != nil {
		if !service.IsNotFoundError(err) {
			c.JSON(500, handleError(err))
			return
		}
	} else {
		c.JSON(400, handleError(fmt.Errorf("user with the same email already exists")))
		return
	}

	newUser := &model.User{
		ID:        uuid.New().String(),
		FirstName: dtoBody.FirstName,
		LastName:  dtoBody.LastName,
		Email:     dtoBody.Email,
		Password:  dtoBody.Password,
	}

	err = h.Router.Account.Register(newUser)
	if err != nil {
		c.JSON(500, handleError(err))
		return
	}

	token, err := h.Router.Jwt.GenToken(newUser.ID, newUser.Salt)
	if err != nil {
		c.JSON(500, handleError(err))
		return
	}
	c.JSON(201, gin.H{"token": token})
}
