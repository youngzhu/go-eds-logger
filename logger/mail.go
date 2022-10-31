package logger

import (
	"crypto/tls"
	gomail "gopkg.in/gomail.v2"
	"log"
	"os"
)

var from, fromPwd, to string

func init() {
	from, _ = os.LookupEnv("MAIL_FROM")
	fromPwd, _ = os.LookupEnv("MAIL_FROM_PWD")
	to, _ = os.LookupEnv("MAIL_TO")
}

func SendMail(subject, body string) {
	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)

	m.SetHeader("Subject", subject)

	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.163.com", 465, from, fromPwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		log.Print("发送邮件失败：", err)
	}

}
