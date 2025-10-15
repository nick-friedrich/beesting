package mail

import "fmt"

type ConsoleAdapter struct{}

func (c *ConsoleAdapter) Send(email *Email) error {
	fmt.Printf(
		"\n📧 [ConsoleMailer]\nFrom: %s\nTo: %s\nSubject: %s\nBody:\n%s\n\n",
		email.From,
		email.To,
		email.Subject,
		email.Body,
	)
	return nil
}
