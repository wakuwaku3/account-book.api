package application

import (
	"strings"
)

type (
	// ErrorCode はエラーの種類を識別します
	ErrorCode string

	clientError struct {
		codes *[]ErrorCode
	}
	// ClientError はクライアントサイドでハンドルするアプリケーション固有のエラーです
	ClientError interface {
		error
		HasError() bool
		Append(codes ...ErrorCode)
		GetErrorCodes() *[]string
	}
)

const (
	// RequiredMailAddress :メールアドレスは必須です
	RequiredMailAddress ErrorCode = "00001"
	// RequiredPassword :パスワードは必須です
	RequiredPassword ErrorCode = "00002"
	// RequiredPasswordToken :パスワードトークンは必須です
	RequiredPasswordToken ErrorCode = "00003"
	// RequiredAmount :金額は必須です
	RequiredAmount ErrorCode = "00004"
	// RequiredCategory :カテゴリは必須です
	RequiredCategory ErrorCode = "00005"
	// ClosedTransaction :取引はClosedです
	ClosedTransaction ErrorCode = "00006"
	// RequiredID :IDは必須です
	RequiredID ErrorCode = "00007"
	// RequiredPlanName :計画名は必須です
	RequiredPlanName ErrorCode = "00008"
	// MoreThanZeroInterval :間隔は0より大きい必要があります
	MoreThanZeroInterval ErrorCode = "00009"
	// InValidDateRange :日付の範囲が不正です
	InValidDateRange ErrorCode = "00010"
	// NotFound :データが存在しません
	NotFound ErrorCode = "00011"
	// FailureSignIn :サインインに失敗しました
	FailureSignIn ErrorCode = "00012"
	// FailurePasswordReset :パスワード変更に失敗しました
	FailurePasswordReset ErrorCode = "00013"
	// IsDeleted :削除済みです
	IsDeleted ErrorCode = "00014"
	// LessLengthPathword :パスワードは8文字以上設定してください。
	LessLengthPathword ErrorCode = "00015"
	// InvalidCharPassword :パスワードには、半角英小文字、大文字、数字をそれぞれ1種類以上使用してください。
	InvalidCharPassword ErrorCode = "00016"
	// RequiredAgreement :サービス利用規約への同意が必要です。
	RequiredAgreement ErrorCode = "00017"
	// FailureSignUp :サインアップに失敗しました。
	FailureSignUp ErrorCode = "00018"
	// ExpiredURL :URLの有効期限が切れています。
	ExpiredURL ErrorCode = "00019"
	// RequiredSignUpToken :サインアップトークンは必須です。
	RequiredSignUpToken ErrorCode = "00020"
	// RequiredName :名前は必須です。
	RequiredName ErrorCode = "00021"
	// InValidCulture :不正なカルチャーです。
	InValidCulture ErrorCode = "00022"
)

// NewClientError はClientErrorを生成します。
func NewClientError(codes ...ErrorCode) ClientError {
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
