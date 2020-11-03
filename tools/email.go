package tools

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

var auth smtp.Auth

func SendHTML(template string, receiver_url string, receiver_name string, url string) {
	auth = smtp.PlainAuth("", "autobotrobobot@gmail.com", "auto123123", "smtp.gmail.com") //In order to auth in gmail, the "allow less secure apps" should be activated. In https://myaccount.google.com/security
	templateData := struct {
		Name string
		URL  string
		Date string
	}{
		Name: receiver_name,
		URL:  url,
		Date: "9 Enero 2007 a las 4:00",
	}
	r := NewRequest([]string{receiver_url}, "Mygame registration", "Su cita")
	err := r.ParseTemplate(template, templateData)
	log.Println(err)
	if err := r.ParseTemplate(template, templateData); err == nil {
		ok, _ := r.SendEmail()
		fmt.Println(ok)
	}

}

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, "autobotrobobot@gmail.com", r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
