package service

import (
	"go-postgres/pkg/model"
	"go-postgres/pkg/repository"
	"go-postgres/pkg/repository/query"
)

type Role struct {
	repository *repository.Role
}

func NewRole(r *repository.Role) *Role {
	return &Role{r}
}

func (s *Role) GetList() ([]*model.Role, error) {
	return s.repository.Find(query.Order{"name": "asc"})
}

func (s *Role) GetListByUserID(userID string) ([]*model.Role, error) {
	return s.repository.FindByUserID(userID)
}

func (s *Role) GetByName(name string) (*model.Role, error) {
	return s.repository.FindOne(query.Where{"name": name})
}
