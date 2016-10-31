package logrus_mail

import "testing"

func TestNewPlainAuthMail(t *testing.T) {

	// Check invalid from address
	_, err := NewPlainAuthMail("127.0.0.1", 25, "metest.com", []string{"you@test.com"}, "me@test.com", "asdasd")
	if err == nil {
		t.Errorf("Invalid from address did not return error")
	}

	// Check invalid recipient address
	_, err2 := NewPlainAuthMail("127.0.0.1", 25, "me@test.com", []string{"you@test.com", "badaddress", "other@test.com"}, "me@test.com", "asdasd")
	if err2 == nil {
		t.Errorf("Invalid recipient did not return error")
	}

}
