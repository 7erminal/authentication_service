package functions

import (
	"net/smtp"

	"github.com/beego/beego/v2/core/logs"
)

func SendEmail(username string, otp string) {
	// create app password in gmail to use here. This is different from your login password. This email will send emails.
	auth := smtp.PlainAuth("", "info@amcrentalsgh.com", "kvfb hrjt qmyr lrzm", "smtp.gmail.com")

	// Here we do it all: connect to our server, set up a message and send it

	to := []string{username}

	msg := []byte("To: " + username + "\r\n" +

		"Subject: Your One time pin. Adepa.\r\n" +

		"\r\n" +

		"Your one time pin is " + otp + ".\r\nThis code will expire in 5 mins.\r\n")

	err := smtp.SendMail("smtp.gmail.com:587", auth, "info@amcrentalsgh.com", to, msg)

	if err != nil {

		logs.Debug(err)

	}
}

func SendEmailNew(email string, subject_ string, message string) {
	// Create app password in gmail to use here
	auth := smtp.PlainAuth("", "info@amcrentalsgh.com", "@Amcadmin2025", "smtp.gmail.com")

	// Here we do it all: connect to our server, set up a message and send it

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject:" + subject_
	body := message

	to := []string{email}

	msg := []byte(subject + mime + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, "info@amcrentalsgh.com", to, msg)

	if err != nil {

		logs.Debug(err)

	}
}
