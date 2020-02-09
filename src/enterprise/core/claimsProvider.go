package core

// ClaimsProvider は Claimsを取得します
type ClaimsProvider interface {
	GetUserID() *string
	GetEmail() *string
	Authenticated() bool
}
