package sendgrid

import (
	"encoding/json"

	sg "github.com/sendgrid/sendgrid-go"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
)

type (
	helper struct {
		env application.Env
	}
	// Helper はSendGridのヘルパークラスです
	Helper interface {
		Send(body *RequestBody) error
	}
	// RequestBody は送信内容です
	RequestBody struct {
		From             MailAddress       `json:"from"`
		Personalizations []Personalization `json:"personalizations,omitempty"`
		TemplateID       string            `json:"template_id"`
	}
	// MailAddress です
	MailAddress struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	}
	// Personalization です
	Personalization struct {
		To                  []MailAddress     `json:"to"`
		DynamicTemplateData map[string]string `json:"dynamic_template_data,omitempty"`
	}
)

// NewHelper is create instance
func NewHelper(env application.Env) Helper {
	return &helper{env}
}
func (t *helper) Send(body *RequestBody) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	apiKey := t.env.GetSendGridAPIKey()
	request := sg.GetRequest(*apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = b
	err = core.Try(func() error {
		_, err := sg.API(request)
		return err
	}, 10)
	return err
}
