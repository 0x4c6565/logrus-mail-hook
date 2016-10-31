# logrus-mail-hook

Usage
-----

```go

// Initialise mail interface implementation NewAuthMail for plain-text authentication
mail, err := logrus_mail.NewPlainAuthMail("127.0.0.1", 25, "me@test.com", []string{"you@test.com"}, "me@test.com", "reallystrongpassword")
if err != nil {
  println(fmt.Sprintf("Failed to initialise mail hook mailer [%s]", err))
}

// Add mail interface to mail hook with error level
mailHook, err := logrus_mail.NewMailHook(mail, []log.Level{log.ErrorLevel})
if err != nil {
  println(fmt.Sprintf("Failed to create mail hook [%s]", err))
}

// Add mail hook to logrus
log.AddHook(mailHook)

// Fire a test
log.Errorf("test")
```
