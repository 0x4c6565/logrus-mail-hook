package logrus_mail

import "github.com/sirupsen/logrus"

type MailHook struct {
	mail   Mail
	levels []logrus.Level
}

// NewMailHook returns pointer to MailHook
func NewMailHook(mail Mail, levels []logrus.Level) (*MailHook, error) {
	return &MailHook{
		mail:   mail,
		levels: levels,
	}, nil
}

// Levels returns the available logging levels.
func (m *MailHook) Levels() []logrus.Level {
	return m.levels
}

// Fire fires entry off to mail
func (m *MailHook) Fire(entry *logrus.Entry) error {
	return m.mail.Send(entry)
}
