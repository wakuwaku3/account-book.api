package application

import "github.com/wakuwaku3/account-book.api/src/enterprise/domains/core"

const (
	// RequiredMailAddress :メールアドレスは必須です
	RequiredMailAddress core.ErrorCode = "00001"
	// RequiredPassword :パスワードは必須です
	RequiredPassword core.ErrorCode = "00002"
	// RequiredPasswordToken :パスワードトークンは必須です
	RequiredPasswordToken core.ErrorCode = "00003"
	// RequiredAmount :金額は必須です
	RequiredAmount core.ErrorCode = "00004"
	// RequiredCategory :カテゴリは必須です
	RequiredCategory core.ErrorCode = "00005"
	// ClosedTransaction :取引はClosedです
	ClosedTransaction core.ErrorCode = "00006"
	// RequiredID :IDは必須です
	RequiredID core.ErrorCode = "00007"
	// RequiredPlanName :計画名は必須です
	RequiredPlanName core.ErrorCode = "00008"
	// MoreThanZeroInterval :間隔は0より大きい必要があります
	MoreThanZeroInterval core.ErrorCode = "00009"
	// InValidDateRange :日付の範囲が不正です
	InValidDateRange core.ErrorCode = "00010"
	// NotFound :データが存在しません
	NotFound core.ErrorCode = "00011"
	// FailureSignIn :サインインに失敗しました
	FailureSignIn core.ErrorCode = "00012"
	// FailurePasswordReset :パスワード変更に失敗しました
	FailurePasswordReset core.ErrorCode = "00013"
	// IsDeleted :削除済みです
	IsDeleted core.ErrorCode = "00014"
	// LessLengthPathword :パスワードは8文字以上設定してください。
	LessLengthPathword core.ErrorCode = "00015"
	// InvalidCharPassword :パスワードには、半角英小文字、大文字、数字をそれぞれ1種類以上使用してください。
	InvalidCharPassword core.ErrorCode = "00016"
	// RequiredAgreement :サービス利用規約への同意が必要です。
	RequiredAgreement core.ErrorCode = "00017"
	// FailureSignUp :サインアップに失敗しました。
	FailureSignUp core.ErrorCode = "00018"
	// ExpiredURL :URLの有効期限が切れています。
	ExpiredURL core.ErrorCode = "00019"
	// RequiredSignUpToken :サインアップトークンは必須です。
	RequiredSignUpToken core.ErrorCode = "00020"
	// RequiredName :名前は必須です。
	RequiredName core.ErrorCode = "00021"
	// InValidCulture :不正なカルチャーです。
	InValidCulture core.ErrorCode = "00022"
)
