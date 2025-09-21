package sendmail

import (
	"net/smtp"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/jordan-wright/email"
)

func SendEmailPasswordReset(config *configs.Config, to, token string) error {
	e := email.NewEmail()
	e.From = config.EmailSender.Email
	e.To = []string{to}
	e.Subject = "Password reset"

	// TODO: Вставить нужную ссылку для перехода к сбросу пароля и созданию нового
	e.HTML = []byte("<p>Please click the following link to reset your password: <a href=\"#" + token + "\">Password reset</a></p>")
	e.Text = []byte("Please click the following link to reset your password: http://example/verify/" + token)

	err := e.Send(
		config.EmailSender.SmtpWithPort,
		smtp.PlainAuth("", config.EmailSender.Email, config.EmailSender.Password, config.EmailSender.Smtp),
	)
	if err != nil {
		return err
	}
	return nil
}
