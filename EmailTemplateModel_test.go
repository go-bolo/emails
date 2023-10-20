package emails

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmailTemplateModel_Render(t *testing.T) {
	type fields struct {
		ID        uint64
		Subject   string
		Text      string
		Css       string
		Html      string
		Type      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	type args struct {
		ctx TemplateVariables
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Email
		wantErr bool
	}{
		{
			name: "success on simple data",
			fields: fields{
				ID:        1,
				Subject:   "Hello {{name}}, this is a test",
				Text:      "Hello {{name}}, this is a test in body field",
				Css:       "",
				Html:      "<p>Hello {{name}}, this is a test in html body field</p>",
				Type:      "type",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			args: args{
				ctx: TemplateVariables{
					"name": "John",
				},
			},

			want: &Email{
				Subject: "Hello John, this is a test",
				Text:    "Hello John, this is a test in body field",
				HTML:    "<p>Hello John, this is a test in html body field</p>",
			},
		},

		{
			name: "success with css",
			fields: fields{
				ID:        1,
				Subject:   "Hello {{name}}, this is a test",
				Text:      "Hello {{name}}, this is a test in body field",
				Css:       "p { backgroud-color: red; }",
				Html:      "<p>Hello {{name}}, this is a test in html body field</p>",
				Type:      "type",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			args: args{
				ctx: TemplateVariables{
					"name": "John",
				},
			},

			want: &Email{
				Subject: "Hello John, this is a test",
				Text:    "Hello John, this is a test in body field",
				HTML:    "<html><head></head><body><p style=\"backgroud-color:red\">Hello John, this is a test in html body field</p></body></html>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EmailTemplateModel{
				ID:        tt.fields.ID,
				Subject:   tt.fields.Subject,
				Text:      tt.fields.Text,
				Css:       tt.fields.Css,
				Html:      tt.fields.Html,
				Type:      tt.fields.Type,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			got := Email{}
			err := r.Render(tt.args.ctx, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailTemplateModel.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, &got)
		})
	}
}
