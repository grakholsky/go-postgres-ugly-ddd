package repository

import (
	"errors"
	"go-postgres/pkg/model"
	"go-postgres/pkg/repository/query"
	"github.com/jinzhu/gorm"
	"reflect"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{db}
}

func (r *User) FindOne(where query.Where) (*model.User, error) {
	var user model.User
	return &user, where.Make(r.db).Take(&user).Error
}

func (r *User) Find(args ...interface{}) ([]*model.User, error) {
	var users []*model.User
	return users, query.Args(r.db, args...).Find(&users).Error
}

func (r *User) Count() (uint, error) {
	var count uint
	return count, r.db.Table("users").Count(&count).Error
}

func (r *User) Register(models ...interface{}) error {
	if len(models) == 0 {
		return errors.New("register models is empty")
	}
	tx := r.db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	for _, v := range models {
		if v != nil {
			reflectValue := reflect.ValueOf(v)
			if reflectValue.Kind() == reflect.Ptr && !reflectValue.IsNil() {
				reflectValue = reflectValue.Elem()
			}
			if reflectValue.Kind() == reflect.Struct {
				if err := tx.Create(v).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	return tx.Commit().Error
}
