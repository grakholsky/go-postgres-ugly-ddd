package service

import "go-postgres/pkg/repository"

func IsNotFoundError(err error) bool {
	return repository.IsNotFoundError(err)
}
