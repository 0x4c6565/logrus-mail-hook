package logrus_mail

import (
	"bytes"
	"fmt"
	"net/mail"
	"net/smtp"
	"regexp"
	"strings"

	"os"

	"github.com/sirupsen/logrus"
)

type Mail interface {
	Send(entry *logrus.Entry) error
}

type PlainAuthMail struct {
	host         string
	port         int
	from         string
	recipients   []string
	authUser     string
	authPassword string
}

var rCRLF *regexp.Regexp

func init() {
	rCRLF = regexp.MustCompile(`\r?\n`)
}

func NewPlainAuthMail(host string, port int, from string, recipients []string, authUser string, authPassword string) (Mail, error) {
	mail := PlainAuthMail{
		host:         host,
		port:         port,
		from:         from,
		recipients:   recipients,
		authUser:     authUser,
		authPassword: authPassword,
	}

	err := mail.validateParameters()
	if err != nil {
		return nil, err
	}

	return mail, nil
}

func (m PlainAuthMail) Send(entry *logrus.Entry) error {
	auth := smtp.PlainAuth(
		"",
		m.authUser,
		m.authPassword,
		m.host,
	)

	message := m.getMessage(entry)

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", m.host, m.port),
		auth,
		m.from,
		m.recipients,
		message,
	)
}

func (m *PlainAuthMail) validateParameters() error {
	if m.host == "" {
		return fmt.Errorf("Valid SMTP host must be provided")
	}
	if m.port < 1 || m.port > 65535 {
		return fmt.Errorf("Valid SMTP port must be provided")
	}

	_, err := mail.ParseAddress(m.from)
	if err != nil {
		return fmt.Errorf("From address failed validation: %s", err)
	}

	_, err = mail.ParseAddress(m.from)
	if err != nil {
		return fmt.Errorf("From address failed validation: %s", err)
	}

	for _, recipient := range m.recipients {
		_, err = mail.ParseAddress(recipient)
		if err != nil {
			return fmt.Errorf("One or more recipients failed validation: %s", err)
		}
	}

	return nil
}

func (m *PlainAuthMail) getMessage(entry *logrus.Entry) []byte {
	entryMessage := formatEntryMessage(entry)
	entryData := formatEntryData(entry)
	subject := getMessageSubject(entry)

	alloc := 4
	headers := make(map[string]string, alloc)
	headers["From"] = m.from
	headers["To"] = strings.Join(m.recipients, ",")
	headers["Subject"] = subject
	headers["Content-Type"] = "text/plain; charset=\"utf-8\""

	buf := bytes.Buffer{}
	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString("Fields:\r\n")
	buf.WriteString(entryData)
	buf.WriteString("Message:\r\n")
	buf.WriteString(entryMessage)
	buf.WriteString("\r\n")

	return buf.Bytes()
}

func formatEntryMessage(entry *logrus.Entry) string {
	return rCRLF.ReplaceAllString(entry.Message, "\r\n")
}

func formatEntryData(entry *logrus.Entry) string {
	buf := bytes.Buffer{}

	for k, v := range entry.Data {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	return buf.String()
}

func getMessageSubject(entry *logrus.Entry) string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("[%s] Log entry on %s", entry.Level.String(), hostname)
}
