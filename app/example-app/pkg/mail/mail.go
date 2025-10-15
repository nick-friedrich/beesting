package mail

import (
	"errors"
	"sync"
)

type Mailer struct {
	Adapter MailerAdapter
}

type MailerAdapter interface {
	Send(email *Email) error
}

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

var (
	mailerInstance *Mailer
	once           sync.Once
)

func GetMailer() *Mailer {
	once.Do(func() {
		mailerInstance = &Mailer{}
	})
	return mailerInstance
}

func InitMailer(adapter MailerAdapter) {
	once.Do(func() {
		mailerInstance = &Mailer{
			Adapter: adapter,
		}
	})
}

func (m *Mailer) SendEmail(email *Email) error {
	if m.Adapter == nil {
		return errors.New("mailer adapter not initialized")
	}
	return m.Adapter.Send(email)
}
