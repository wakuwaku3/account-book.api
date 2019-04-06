package mails

import (
	"encoding/json"
	"net/url"
	"path"

	"github.com/wakuwaku3/account-book.api/src/application"
)

type (
	userExisting struct {
		env    application.Env
		helper Helper
	}
)

// NewUserExisting is create instance
func NewUserExisting(env application.Env, helper Helper) application.UserExistingMail {
	return &userExisting{env, helper}
}
func (t *userExisting) Send(args *application.UserExistingMailSendArgs) error {
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
		TemplateID: "d-11ea8314916e4ccdb0ac934ce87efb8f",
	}
	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	return t.helper.Send(&body)
}
