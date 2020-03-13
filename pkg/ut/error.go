package ut

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

/**
快速检查错误是否是数字转化的错误
*/
func IsInvalidPager(err error) bool {
	if err == nil {
		return false
	}
	isInvalid := func(param string) bool {
		if strings.HasPrefix(err.Error(), fmt.Sprintf("param=%s value=", param)) && strings.Contains(err.Error(), "is invalid: strconv.Atoi: parsing") {
			return true
		}
		return false
	}
	if isInvalid("limit") || isInvalid("page") {
		return true
	}
	return false
}

//数据库事务处理返回的Err 处理是否需要直接返回给用户
func IsUserErr(err error, errors []error) bool {
	for _, userErr := range errors {
		if err == userErr {
			return true
		}
	}
	return false
}

/**
验证错误类，主要用于在验证Form时， 是传入的form数据不合法，还是数据库等其它的错误
*/
type ValidateError struct {
	s string
}

func (err ValidateError) Error() string {
	return err.s
}

func NewValidateError(format string, a ...interface{}) ValidateError {
	return ValidateError{s: fmt.Sprintf(format, a...)}
}

func IsValidateError(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(ValidateError); ok {
		return true
	} else {
		return false
	}
}

type AuthError struct {
	s string
}

func (err AuthError) Error() string {
	return err.s
}

func NewAuthError(format string, a ...interface{}) AuthError {
	return AuthError{s: fmt.Sprintf(format, a...)}
}

func IsAuthError(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(AuthError); ok {
		return true
	} else {
		return false
	}
}

func NewError(format string, a ...interface{}) error {
	return errors.WithStack(errors.New(fmt.Sprintf(format, a...)))
}

type BadCodeError struct {
	InError error
	Code    int
	ErrMsg  string
}

func (err BadCodeError) Error() string {
	if err.InError != nil {
		return err.InError.Error()
	}
	return fmt.Sprintf(`code: %d`, err.Code)
}

func NewBadCodeError(inError error, code int, msg string) BadCodeError {
	return BadCodeError{InError: inError, Code: code, ErrMsg: msg}
}

func IsBadCodeError(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(BadCodeError); ok {
		return true
	} else {
		return false
	}
}
