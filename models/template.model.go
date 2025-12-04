package models

type Template struct {
	Id string `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey" faker:"uuid_hyphenated"`
}

type TemplateCreate struct {
}

type TemplateQuery struct {
	Id *string `form:"id"`
}

func (Template) TableName() string {
	return "templates"
}
