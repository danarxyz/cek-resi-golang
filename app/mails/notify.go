package mails

import (
	"fmt"
	"goravel/app/models"

	"github.com/goravel/framework/contracts/mail"
)

type Notify struct {
	models.Resi
}

func NewNotify(resi models.Resi) *Notify {
	return &Notify{Resi: resi}
}

// Attachments attach files to the mail
func (receiver *Notify) Attachments() []string {
	return []string{}
}

// Content set the content of the mail
func (receiver *Notify) Content() *mail.Content {
	view := fmt.Sprintf(`
        <html>
        <body>
			<h1>Cek Paket Anda</h1>
			<p>Package Name: %s</p>
            <p>Resi: %s</p>
            <p>Detail: %s</p>
        </body>
        </html>
    `, receiver.PackageName, receiver.Resi.TrackingNum, receiver.Resi.Details)
	return &mail.Content{Html: view}
}

// Envelope set the envelope of the mail
func (receiver *Notify) Envelope() *mail.Envelope {
	return &mail.Envelope{
		From:    mail.Address{Address: "noreply@cekresi.com", Name: "Cek Resi"},
		Subject: "Cek Paket Anda",
		To:      []string{receiver.Resi.Email},
	}
}

// Queue set the queue of the mail
func (receiver *Notify) Queue() *mail.Queue {
	return &mail.Queue{}
}
