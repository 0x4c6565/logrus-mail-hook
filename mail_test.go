package logrus_mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlainAuthMail(t *testing.T) {
	t.Run("InvalidHost_ReturnsError", func(t *testing.T) {
		_, err := NewPlainAuthMail("", 25, "me@test.com", []string{"you@test.com"}, "auth@test.com", "testpassword")

		assert.NotNil(t, err)
		assert.Equal(t, "Valid SMTP host must be provided", err.Error())
	})

	t.Run("InvalidPort_ReturnsError", func(t *testing.T) {
		_, err := NewPlainAuthMail("127.0.0.1", 0, "me@test.com", []string{"you@test.com"}, "auth@test.com", "testpassword")

		assert.NotNil(t, err)
		assert.Equal(t, "Valid SMTP port must be provided", err.Error())
	})

	t.Run("InvalidFromAddress_ReturnsError", func(t *testing.T) {
		_, err := NewPlainAuthMail("127.0.0.1", 25, "metest.com", []string{"you@test.com"}, "auth@test.com", "testpassword")

		assert.NotNil(t, err)
		assert.Equal(t, "From address failed validation: mail: no angle-addr", err.Error())
	})

	t.Run("InvalidRecipientAddress_ReturnsError", func(t *testing.T) {
		_, err := NewPlainAuthMail("127.0.0.1", 25, "me@test.com", []string{"you@test.com", "badaddress", "other@test.com"}, "auth@test.com", "asdasd")

		assert.NotNil(t, err)
		assert.Equal(t, "One or more recipients failed validation: mail: no angle-addr", err.Error())
	})
}
