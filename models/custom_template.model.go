package models

type CustomTemplate struct {
}

type CustomTemplateCreate struct {
	CustomTemplate
}

type CustomTemplatePatch struct {
	CustomTemplate
}

type CustomTemplateQuery struct {
	CustomTemplate
}

type CustomTemplateResponse struct {
}

func (CustomTemplate) TableName() string {
	return "templates"
}
