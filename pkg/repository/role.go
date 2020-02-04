package repository

import (
	"go-postgres/pkg/model"
	"go-postgres/pkg/repository/query"
	"github.com/jinzhu/gorm"
)

type Role struct {
	db *gorm.DB
}

func NewRole(db *gorm.DB) *Role {
	return &Role{db}
}

func (r *Role) Find(args ...interface{}) ([]*model.Role, error) {
	var roles []*model.Role
	return roles, query.Args(r.db, args...).Find(&roles).Error
}

func (r *Role) FindOne(where query.Where) (*model.Role, error) {
	var role model.Role
	return &role, where.Make(r.db).Take(&role).Error
}

func (r *Role) FindByUserID(userID string) ([]*model.Role, error) {
	var roles []*model.Role
	return roles, r.db.Table(
		"user_roles",
	).Select(
		`roles.id, roles.name`,
	).Joins(
		"left join roles on roles.id = user_roles.role_id",
	).Where(
		"user_id = ?", userID,
	).Scan(&roles).Error
}
