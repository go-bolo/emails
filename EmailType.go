package emails

type EmailType struct {
	Label             string                       `json:"label"`
	DefaultSubject    string                       `json:"defaultSubject"`
	DefaultHTML       string                       `json:"defaultHTML"`
	DefaultText       string                       `json:"defaultText"`
	TemplateVariables map[string]*TemplateVariable `json:"templateVariables"`
}

type TemplateVariable struct {
	Example     string `json:"example"`
	Description string `json:"description"`
}

type EmailTypes map[string]*EmailType
