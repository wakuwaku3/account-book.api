package mails

import (
	"net/url"
	"path"

	"github.com/wakuwaku3/account-book.api/src/adapter/mails/sendgrid"

	"github.com/wakuwaku3/account-book.api/src/application"
)

type (
	resetPassword struct {
		env    application.Env
		helper sendgrid.Helper
	}
)

// NewResetPassword is create instance
func NewResetPassword(env application.Env, helper sendgrid.Helper) application.ResetPasswordMail {
	return &resetPassword{env, helper}
}
func (t *resetPassword) Send(args *application.ResetPasswordMailSendArgs) error {
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
		TemplateID: "d-c9a6a3360fd14791bc63eb3cb3682e90",
	}
	return t.helper.Send(b)
}
