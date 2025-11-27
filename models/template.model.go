package models

type Template struct {
	Id string `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey" faker:"uuid_hyphenated"`
}

func (Template) TableName() string {
	return "templates"
}
