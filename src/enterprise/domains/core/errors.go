package core

import (
	"strings"
)

type (
	// ErrorCode はエラーの種類を識別します
	ErrorCode string

	clientError struct {
		codes *[]ErrorCode
	}
	// Error はクライアントサイドでハンドルするアプリケーション固有のエラーです
	Error interface {
		error
		HasError() bool
		Append(codes ...ErrorCode)
		Concat(Error)
		getErrorCodes() *[]ErrorCode
		GetErrorCodes() *[]string
	}
)

const (
	// RequiredID :IDは必須です
	RequiredID ErrorCode = "core-00001"
	// NotFound :データが存在しません
	NotFound ErrorCode = "core-00002"
)

// NewError はErrorを生成します。
func NewError(codes ...ErrorCode) Error {
	return &clientError{codes: &codes}
}

func (t *clientError) Error() string {
	array := *t.GetErrorCodes()
	return "an error has occurred [" + strings.Join(array, ",") + "]"
}

func (t *clientError) Append(codes ...ErrorCode) {
	array := *t.codes
	for _, code := range codes {
		array = append(array, code)
	}
	t.codes = &array
}

func (t *clientError) Concat(ce Error) {
	codes := ce.getErrorCodes()
	array := *codes
	t.Append(array...)
}

func (t *clientError) getErrorCodes() *[]ErrorCode {
	return t.codes
}

func (t *clientError) GetErrorCodes() *[]string {
	array := make([]string, len(*t.codes))
	for i, code := range *t.codes {
		array[i] = string(code)
	}
	return &array
}

func (t *clientError) HasError() bool {
	return len(*t.codes) > 0
}
