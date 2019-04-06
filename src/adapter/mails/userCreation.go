package mails

import (
	"net/url"
	"path"

	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/drivers/sendgrid"
)

type (
	userCreation struct {
		env    application.Env
		helper sendgrid.Helper
	}
)

// NewUserCreation is create instance
func NewUserCreation(env application.Env, helper sendgrid.Helper) application.UserCreationMail {
	return &userCreation{env, helper}
}
func (t *userCreation) Send(args *application.UserCreationMailSendArgs) error {
	u, _ := url.Parse(*t.env.GetFrontEndURL())
	u.Path = path.Join("sign-up", args.Token)
	b := &sendgrid.RequestBody{
		From: sendgrid.MailAddress{
			Name:  "Account Book Support",
			Email: "support@prj-account-book.firebaseapp.com",
		},
		Personalizations: []sendgrid.Personalization{
			sendgrid.Personalization{
				To: []sendgrid.MailAddress{
					sendgrid.MailAddress{
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
	return t.helper.Send(b)
}
