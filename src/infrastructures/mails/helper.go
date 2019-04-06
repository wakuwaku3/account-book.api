package mails

import (
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/helpers"
)

type (
	helper struct {
		env application.Env
	}
	// Helper はSendGridのヘルパークラスです
	Helper interface {
		Send(body *[]byte) error
	}
	requestBody struct {
		From             mailAddress       `json:"from"`
		Personalizations []personalization `json:"personalizations,omitempty"`
		TemplateID       string            `json:"template_id"`
	}
	mailAddress struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	}
	personalization struct {
		To                  []mailAddress     `json:"to"`
		DynamicTemplateData map[string]string `json:"dynamic_template_data,omitempty"`
	}
)

// NewHelper is create instance
func NewHelper(env application.Env) Helper {
	return &helper{env}
}
func (t *helper) Send(body *[]byte) error {
	apiKey := t.env.GetSendGridAPIKey()
	request := sendgrid.GetRequest(*apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = *body
	err := helpers.Try(func() error {
		_, err := sendgrid.API(request)
		return err
	}, 10)
	return err
}
