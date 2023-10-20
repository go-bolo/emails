package emails

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func SendEmailAsync(opts *EmailOpts) {
	go func() {
		email, err := NewEmailWithTemplate(opts)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": fmt.Sprintf("%+v\n", err),
			}).Error("emails.SendEmailAsync error on NewEmailWithTemplate")
		}

		err = email.Send()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": fmt.Sprintf("%+v\n", err),
			}).Error("emails.SendEmailAsync error on email.Send")
		}
	}()
}
