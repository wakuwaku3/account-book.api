package mails

import (
	"encoding/json"
	"net/url"
	"path"

	"github.com/wakuwaku3/account-book.api/src/application"
)

type (
	resetPassword struct {
		env    application.Env
		helper Helper
	}
)

// NewResetPassword is create instance
func NewResetPassword(env application.Env, helper Helper) application.ResetPasswordMail {
	return &resetPassword{env, helper}
}
func (t *resetPassword) Send(args *application.ResetPasswordMailSendArgs) error {
	u, _ := url.Parse(*t.env.GetFrontEndURL())
	u.Path = path.Join("reset-password", args.Token)
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
		TemplateID: "d-c9a6a3360fd14791bc63eb3cb3682e90",
	}
	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	return t.helper.Send(&body)
}
