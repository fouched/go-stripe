package main

import (
	"bytes"
	"embed"
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"time"
)

// with the go embed directive below we can compile
// the templates with the application in a single binary
//
//go:embed templates
var emailTemplateFS embed.FS

func (app *application) SendMail(from, to, subject, tmpl string, data interface{}) error {
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)
	t, err := template.New("email-html").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	var tplHTML bytes.Buffer
	if err = t.ExecuteTemplate(&tplHTML, "body", data); err != nil {
		app.errorLog.Println(err)
		return err
	}

	formattedMessage := tplHTML.String()

	var tplPlain bytes.Buffer
	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	t, err = template.New("email-plain").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	if err = t.ExecuteTemplate(&tplPlain, "body", data); err != nil {
		app.errorLog.Println(err)
		return err
	}

	plainMessage := tplPlain.String()

	app.infoLog.Println(formattedMessage, plainMessage)

	// send the mail
	server := mail.NewSMTPClient()
	server.Host = app.config.smtp.host
	server.Port = app.config.smtp.port
	server.Username = app.config.smtp.username
	server.Password = app.config.smtp.password
	server.Encryption = getEncryption(app.config.smtp.encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, formattedMessage)
	email.AddAlternative(mail.TextPlain, plainMessage)

	err = email.Send(smtpClient)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	app.infoLog.Println("sent mail")

	return nil
}

func getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
