package dto

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"mime/multipart"
	"reflect"
	"strconv"
)

//  DTO (Data Transfer Object) schema

const (
	BindFORM = "form"
	BindJSON = "json"
	BindURI  = "uri"

	fieldErrorRequired = "%s - this field is required"
	fieldErrorEmail    = "%s - this field has wrong email format"
	fieldErrorUUID4    = "%s - this field has wrong uuid4 format; use instead: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	fieldErrorNum      = "%s has wrong format or %v"
)

type (
	Query struct {
		Page     uint   `form:"page"`
		PageSize uint   `form:"page_size"`
		Search   string `form:"search"`
	}

	ParamID struct {
		ID string `uri:"id" binding:"required,uuid4"`
	}

	SignInUser struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	SignUpUser struct {
		FirstName string                `form:"first_name" binding:"required"`
		LastName  string                `form:"last_name" binding:"required"`
		Email     string                `form:"email" binding:"required,email"`
		Password  string                `form:"password" binding:"required"`
		Avatar    *multipart.FileHeader `form:"avatar"`
	}
)

func Errors(dtoStruct interface{}, err error, bindTag string) []error {
	var errs []error
	switch err := err.(type) {
	case validator.ValidationErrors:
		t := reflect.ValueOf(dtoStruct).Type()
		for _, fieldError := range err {
			structField, _ := t.FieldByName(fieldError.StructField())
			fieldName := structField.Tag.Get(bindTag)
			switch fieldError.Tag() {
			case "required":
				errs = append(errs, fmt.Errorf(fieldErrorRequired, fieldName))
			case "email":
				errs = append(errs, fmt.Errorf(fieldErrorEmail, fieldName))
			case "uuid4":
				errs = append(errs, fmt.Errorf(fieldErrorUUID4, fieldName))
			}
		}
	case *strconv.NumError:
		errs = append(errs, fmt.Errorf(fieldErrorNum, err.Num, err.Err))
	default:
		errs = append(errs, err)
	}
	return errs
}
