package apperrors

import (
	"strings"
)

type (
	// Code はエラーの種類を識別します
	Code string

	clientError struct {
		codes *[]Code
	}
	// ClientError はクライアントサイドでハンドルするアプリケーション固有のエラーです
	ClientError interface {
		error
		HasError() bool
		Append(codes ...Code)
		GetErrorCodes() *[]string
	}
)

const (
	// RequiredMailAddress :メールアドレスは必須です
	RequiredMailAddress Code = "00001"
	// RequiredPassword :パスワードは必須です
	RequiredPassword Code = "00002"
	// RequiredPasswordToken :パスワードトークンは必須です
	RequiredPasswordToken Code = "00003"
	// RequiredAmount :金額は必須です
	RequiredAmount Code = "00004"
	// RequiredCategory :カテゴリは必須です
	RequiredCategory Code = "00005"
	// ClosedTransaction :取引はClosedです
	ClosedTransaction Code = "00006"
)

// NewClientError はClientErrorを生成します。
func NewClientError(codes ...Code) ClientError {
	return &clientError{codes: &codes}
}

func (t *clientError) Error() string {
	array := *t.GetErrorCodes()
	return "an error has occurred [" + strings.Join(array, ",") + "]"
}

func (t *clientError) Append(codes ...Code) {
	array := *t.codes
	for _, code := range codes {
		array = append(array, code)
	}
	t.codes = &array
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
