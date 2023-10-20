package emails_test

import (
	"reflect"
	"testing"

	emails "github.com/go-bolo/emails"
	"github.com/stretchr/testify/assert"
)

func TestEmail_QueueToSend(t *testing.T) {
	type fields struct {
		Identifier       string
		Template         *string
		To               string
		From             string
		ReplyTo          string
		CC               string
		CCO              string
		Subject          string
		Text             string
		HTML             string
		DeliveryAttempts int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &emails.Email{
				Identifier:       tt.fields.Identifier,
				Template:         tt.fields.Template,
				To:               tt.fields.To,
				From:             tt.fields.From,
				ReplyTo:          tt.fields.ReplyTo,
				CC:               tt.fields.CC,
				CCO:              tt.fields.CCO,
				Subject:          tt.fields.Subject,
				Text:             tt.fields.Text,
				HTML:             tt.fields.HTML,
				DeliveryAttempts: tt.fields.DeliveryAttempts,
			}
			if err := r.QueueToSend(); (err != nil) != tt.wantErr {
				t.Errorf("Email.QueueToSend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmail_Requeue(t *testing.T) {
	type fields struct {
		Identifier       string
		Template         *string
		To               string
		From             string
		ReplyTo          string
		CC               string
		CCO              string
		Subject          string
		Text             string
		HTML             string
		DeliveryAttempts int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &emails.Email{
				Identifier:       tt.fields.Identifier,
				Template:         tt.fields.Template,
				To:               tt.fields.To,
				From:             tt.fields.From,
				ReplyTo:          tt.fields.ReplyTo,
				CC:               tt.fields.CC,
				CCO:              tt.fields.CCO,
				Subject:          tt.fields.Subject,
				Text:             tt.fields.Text,
				HTML:             tt.fields.HTML,
				DeliveryAttempts: tt.fields.DeliveryAttempts,
			}
			if err := r.Requeue(); (err != nil) != tt.wantErr {
				t.Errorf("Email.Requeue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEmail(t *testing.T) {
	type args struct {
		opts     *emails.EmailOpts
		template *emails.EmailTemplateModel
	}
	tests := []struct {
		name    string
		args    args
		want    *emails.Email
		wantErr bool
	}{
		{
			name: "success on happy path",
			args: args{
				template: &emails.EmailTemplateModel{
					Subject: "Hello {{name}}, this is a test. Your age is {{age}}",
					Text:    "Hello {{name}}, this is a test in body field. Your age is {{age}}",
					Css:     "",
					Html:    "<p>Hello {{name}}, this is a test in html body field.</p> <p>Your age is {{age}}</p>",
					Type:    "test",
				},
				opts: &emails.EmailOpts{
					TemplateName: "test",
					Variables: emails.TemplateVariables{
						"name": "Luffy",
						"age":  "19",
					},
				},
			},
			want: &emails.Email{
				Subject: "Hello Luffy, this is a test. Your age is 19",
				Text:    "Hello Luffy, this is a test in body field. Your age is 19",
				HTML:    "<p>Hello Luffy, this is a test in html body field.</p> <p>Your age is 19</p>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := tt.args.template
			err := template.Save()
			assert.Nil(t, err)

			got, err := emails.NewEmailWithTemplate(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmailWithTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got, tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmail() = %v, want %v", got, tt.want)
			}

			err = template.Delete()
			assert.Nil(t, err)
		})
	}
}
