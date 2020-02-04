package handler

import (
	"github.com/gin-gonic/gin"
	"go-postgres/pkg/dto"
)

func handleQuery(c *gin.Context) (*dto.Query, error) {
	var dtoQuery dto.Query
	if err := c.ShouldBindQuery(&dtoQuery); err != nil {
		return nil, err
	}
	if dtoQuery.Page == 0 {
		dtoQuery.Page = 1
	}
	if dtoQuery.PageSize == 0 {
		dtoQuery.PageSize = 100
	}
	return &dto.Query{
		Page:     (dtoQuery.Page - 1) * dtoQuery.PageSize,
		PageSize: dtoQuery.PageSize,
		Search:   dtoQuery.Search,
	}, nil
}
