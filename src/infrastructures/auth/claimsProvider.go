package auth

import "github.com/wakuwaku3/account-book.api/src/domains"

type (
	claimsProvider struct {
		email  string
		userID string
	}
)

// NewClaimsProvider is create instance
func NewClaimsProvider(email string, userID string) domains.ClaimsProvider {
	return &claimsProvider{email, userID}
}

// NewAnonymousClaimsProvider is create instance
func NewAnonymousClaimsProvider() domains.ClaimsProvider {
	return &claimsProvider{email: "anonymous@example.com", userID: "anonymous"}
}
func (t *claimsProvider) GetEmail() *string {
	return &t.email
}
func (t *claimsProvider) GetUserID() *string {
	return &t.userID
}
