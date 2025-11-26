package models

type Template struct {
	Id string `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
}

func (Template) TableName() string {
	return "templates"
}
