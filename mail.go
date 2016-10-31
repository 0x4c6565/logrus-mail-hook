package logrus_mail

import (
	"bytes"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"time"

	"os"

	"github.com/Sirupsen/logrus"
)

type Mail interface {
	Send(entry *logrus.Entry) error
}

type AuthMail struct {
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

func NewMail(host string, port int, from string, recipients []string, authUser string, authPassword string) (Mail, error) {
	mail := AuthMail{
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

func (m AuthMail) Send(entry *logrus.Entry) error {
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

func (m *AuthMail) validateParameters() error {
	// Check if server listens on that port.
	conn, err := net.DialTimeout("tcp", m.host+":"+strconv.Itoa(m.port), 3*time.Second)
	if err != nil {
		return fmt.Errorf("Failed to connect to SMTP server [%s]", err)
	}
	defer conn.Close()

	_, err = mail.ParseAddress(m.from)
	if err != nil {
		return fmt.Errorf("From address failed validation [%s]", err)
	}

	for _, recipient := range m.recipients {
		_, err = mail.ParseAddress(recipient)
		if err != nil {
			return fmt.Errorf("One or more recipients failed validation [%s]", err)
		}
	}

	return nil
}

func (m *AuthMail) getMessage(entry *logrus.Entry) []byte {
	message := getMessageData(entry)
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
	buf.WriteString("\r\n")
	buf.WriteString(message)
	buf.WriteString("\r\n")

	return buf.Bytes()
}

func getMessageData(entry *logrus.Entry) string {
	return rCRLF.ReplaceAllString(entry.Message, "\r\n")
}

func getMessageSubject(entry *logrus.Entry) string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("[%s] Log entry on %s", entry.Level.String(), hostname)
}
