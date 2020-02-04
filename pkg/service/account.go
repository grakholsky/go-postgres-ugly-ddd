package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"go-postgres/pkg/model"
	"go-postgres/pkg/repository"
	"go-postgres/pkg/repository/query"
	"golang.org/x/crypto/pbkdf2"
)

type Account struct {
	repoUser *repository.User
	repoRole *repository.Role
}

func NewAccount(user *repository.User, role *repository.Role) *Account {
	return &Account{user, role}
}

func (s *Account) GetByID(userID string) (*model.User, error) {
	return s.repoUser.FindOne(query.Where{"id": userID})
}

func (s *Account) GetByEmail(email string) (*model.User, error) {
	return s.repoUser.FindOne(query.Where{"email": email})
}

func (s *Account) GetList(offset, limit uint, search string) ([]*model.User, error) {
	return s.repoUser.Find(
		query.Offset(offset),
		query.Limit(limit),
		query.Search{Text: search, Fields: []string{"first_name", "last_name", "email"}},
		query.Order{"created_at": "desc"},
	)
}

func (s *Account) GetCount() (uint, error) {
	return s.repoUser.Count()
}

func (s *Account) Register(user *model.User) error {
	accountRole, err := s.repoRole.FindOne(query.Where{"name": repository.AccountRole})
	if err != nil {
		return nil
	}
	role := &model.UserRole{
		ID:     uuid.New().String(),
		UserID: user.ID,
		RoleID: accountRole.ID,
	}
	user.Salt = s.genSalt()
	user.Password = s.genPBKDF2Key(user.Password, user.Salt)
	return s.repoUser.Register(user, role)
}

func (s *Account) CheckPassword(userID, password string) (bool, error) {
	dbUser, err := s.GetByID(userID)
	if err != nil {
		return false, fmt.Errorf("check password failed: %v", err)
	}
	return dbUser.Password == s.genPBKDF2Key(password, dbUser.Salt), nil
}

func (s *Account) genPBKDF2Key(password string, salt []byte) string {
	hash := pbkdf2.Key([]byte(password), salt, 4096, 32, sha1.New)
	return hex.EncodeToString(hash)
}

func (s *Account) genSalt() []byte {
	salt := make([]byte, 128)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	return salt
}
