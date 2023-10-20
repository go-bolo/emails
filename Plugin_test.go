package emails_test

import (
	"reflect"
	"testing"

	"github.com/go-bolo/bolo"
	emails "github.com/go-bolo/emails"
)

func TestEmailPlugin_GetName(t *testing.T) {
	type fields struct {
		Name       string
		EmailTypes map[string]*emails.EmailType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &emails.EmailPlugin{
				Name:       tt.fields.Name,
				EmailTypes: tt.fields.EmailTypes,
			}
			if got := p.GetName(); got != tt.want {
				t.Errorf("EmailPlugin.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailPlugin_Init(t *testing.T) {

	type fields struct {
		Name       string
		EmailTypes map[string]*emails.EmailType
	}
	type args struct {
		app bolo.App
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &emails.EmailPlugin{
				Name:       tt.fields.Name,
				EmailTypes: tt.fields.EmailTypes,
			}
			if err := p.Init(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("EmailPlugin.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmailPlugin_BindRoutes(t *testing.T) {
	type fields struct {
		Name       string
		EmailTypes map[string]*emails.EmailType
	}
	type args struct {
		app bolo.App
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &emails.EmailPlugin{
				Name:       tt.fields.Name,
				EmailTypes: tt.fields.EmailTypes,
			}
			if err := p.BindRoutes(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("EmailPlugin.BindRoutes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmailPlugin_AddEmailTemplate(t *testing.T) {
	type fields struct {
		Name       string
		EmailTypes map[string]*emails.EmailType
	}
	type args struct {
		name string
		t    *emails.EmailType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &emails.EmailPlugin{
				Name:       tt.fields.Name,
				EmailTypes: tt.fields.EmailTypes,
			}
			if err := p.AddEmailTemplate(tt.args.name, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("EmailPlugin.AddEmailTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}

			// p.EmailTypes
		})
	}
}

func TestNewPlugin(t *testing.T) {
	type args struct {
		cfg *emails.PluginCfg
	}
	tests := []struct {
		name string
		args args
		want *emails.EmailPlugin
	}{
		{
			name: "TestNewPlugin",
			args: args{
				cfg: &emails.PluginCfg{},
			},
			want: &emails.EmailPlugin{
				Name:       "emails",
				EmailTypes: emails.EmailTypes{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := emails.NewPlugin(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlugin() = %v, want %v", got, tt.want)
			}
		})
	}
}
