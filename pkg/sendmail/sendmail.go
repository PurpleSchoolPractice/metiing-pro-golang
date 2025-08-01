package sendmail

import (
	"net/smtp"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/jordan-wright/email"
)

func SendMail(config *configs.Config, accept, decline string) error {
	e := email.NewEmail()
	e.From = config.EmailSender.Email
	e.To = []string{"tes@test.ru"}

	e.Subject = "Awesome Subject"
	body := decline + " " + accept
	e.Text = []byte(body)

	err := e.Send(config.EmailSender.SmtpWithPort, smtp.PlainAuth("", config.EmailSender.Email, config.EmailSender.Password, config.EmailSender.Smtp))
	if err != nil {
		return err
	}
	return nil
}
