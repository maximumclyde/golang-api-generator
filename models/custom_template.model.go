package models

type CustomTemplate struct {
}

func (CustomTemplate) TableName() string {
	return "templates"
}
