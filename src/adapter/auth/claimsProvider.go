package auth

import "github.com/wakuwaku3/account-book.api/src/application"

type (
	claimsProvider struct {
		email         string
		userID        string
		authenticated bool
	}
)

// NewClaimsProvider is create instance
func NewClaimsProvider(email string, userID string, authenticated bool) application.ClaimsProvider {
	return &claimsProvider{email, userID, authenticated}
}

// NewAnonymousClaimsProvider is create instance
func NewAnonymousClaimsProvider() application.ClaimsProvider {
	return &claimsProvider{email: "anonymous@example.com", userID: "anonymous", authenticated: false}
}
func (t *claimsProvider) GetEmail() *string {
	return &t.email
}
func (t *claimsProvider) GetUserID() *string {
	return &t.userID
}
func (t *claimsProvider) Authenticated() bool {
	return t.authenticated
}
