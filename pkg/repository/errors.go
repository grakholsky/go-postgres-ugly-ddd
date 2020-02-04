package repository

import (
	"github.com/jinzhu/gorm"
)

func IsNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
