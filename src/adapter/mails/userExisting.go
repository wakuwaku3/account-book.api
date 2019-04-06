package mails

import (
	"net/url"
	"path"

	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/sendgrid"
)

type (
	userExisting struct {
		env    application.Env
		helper sendgrid.Helper
	}
)

// NewUserExisting is create instance
func NewUserExisting(env application.Env, helper sendgrid.Helper) application.UserExistingMail {
	return &userExisting{env, helper}
}
func (t *userExisting) Send(args *application.UserExistingMailSendArgs) error {
	u, _ := url.Parse(*t.env.GetFrontEndURL())
	u.Path = path.Join("reset-password", args.Token)
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
		TemplateID: "d-11ea8314916e4ccdb0ac934ce87efb8f",
	}
	return t.helper.Send(b)
}
