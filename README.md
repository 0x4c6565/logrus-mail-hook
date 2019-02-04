# logrus-mail-hook

Usage
-----

```go
package main

import (
	"fmt"

	"github.com/Lee303/logrus-mail-hook"
	log "github.com/Sirupsen/logrus"
)


// Initialise mail interface implementation NewPlainAuthMail for plain-text authentication
mail, err := logrus_mail.NewPlainAuthMail("127.0.0.1", 25, "me@test.com", []string{"you@test.com"}, "me@test.com", "reallystrongpassword")
if err != nil {
  println(fmt.Sprintf("Failed to initialise mail hook mailer: %s", err.Error()))
}

// Add mail interface to mail hook with error level
mailHook, err := logrus_mail.NewMailHook(mail, []log.Level{log.ErrorLevel})
if err != nil {
  println(fmt.Sprintf("Failed to create mail hooks: %s", err.Error()))
}

// Add mail hook to logrus
log.AddHook(mailHook)

// Fire a test
log.Errorf("test")
```
