package logger

import (
	"crypto/tls"
	gomail "gopkg.in/gomail.v2"
	"log"
	"os"
)

func getSecret(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("%s not set\n", key)
	}
	return val
}

func SendMail(subject, body string) {
	from := getSecret("MAIL_FROM")
	fromPwd := getSecret("MAIL_FROM_PWD")
	to := getSecret("MAIL_TO")

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
