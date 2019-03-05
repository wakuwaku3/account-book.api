package mails

import (
	"net/url"
	"path"

	"github.com/wakuwaku3/account-book.api/src/domains"
)

type (
	resetPassword struct {
		env    domains.Env
		helper Helper
	}
)

// NewResetPassword is create instance
func NewResetPassword(env domains.Env, helper Helper) domains.ResetPasswordMail {
	return &resetPassword{env, helper}
}
func (t *resetPassword) Send(args *domains.ResetPasswordMailSendArgs) error {
	u, _ := url.Parse(*t.env.GetFrontEndURL())
	u.Path = path.Join("reset-password", args.Token)

	body := []byte(` {
		"from": {
			"email": "support@prj-account-book.firebaseapp.com",
			"name": "Account Book Support"
		},    
		"personalizations": [
		  {
			"to": [
				{
					"email": "` + args.Email + `",
				}
			],
			"dynamic_template_data":{  
			  "url":"` + u.String() + `",
			}
		  }
		],
		"template_id":"d-c9a6a3360fd14791bc63eb3cb3682e90"
	  }`)

	return t.helper.Send(&body)
}
