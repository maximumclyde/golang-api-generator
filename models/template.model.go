package models

type Template struct {
	Id string `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey" faker:"uuid_hyphenated"`
}

type TemplateCreate struct {
	Template
}

type TemplatePatch struct {
	Template
}

type TemplateQuery struct {
	Template
	Id *string `form:"id"`
}

type TemplateResponse struct {
	Id string `json:"id" faker:"uuid_hyphenated"`
}

func (Template) TableName() string {
	return "templates"
}
