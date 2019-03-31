package mails

import (
	"encoding/json"
	"net/url"
	"path"

	"github.com/wakuwaku3/account-book.api/src/domains"
)

type (
	userCreation struct {
		env    domains.Env
		helper Helper
	}
)

// NewUserCreation is create instance
func NewUserCreation(env domains.Env, helper Helper) domains.UserCreationMail {
	return &userCreation{env, helper}
}
func (t *userCreation) Send(args *domains.UserCreationMailSendArgs) error {
	u, _ := url.Parse(*t.env.GetFrontEndURL())
	u.Path = path.Join("sign-up", args.Token)
	b := &requestBody{
		From: mailAddress{
			Name:  "Account Book Support",
			Email: "support@prj-account-book.firebaseapp.com",
		},
		Personalizations: []personalization{
			personalization{
				To: []mailAddress{
					mailAddress{
						Email: args.Email,
					},
				},
				DynamicTemplateData: map[string]string{
					"url": u.String(),
				},
			},
		},
		TemplateID: "d-2a17843b78824d62835039e74ba0429f",
	}
	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	return t.helper.Send(&body)
}
