package objects

import (
	"project-survey-generator/internal/enums"
)

type Appearance struct {
	Id         int
	Type       enums.EnumAppearanceType
	TemplateId int
	Params     map[string]string
}

func NewAppearance(id int, aType enums.EnumAppearanceType, templateId int, params map[string]string) *Appearance {
	return &Appearance{
		Id:         id,
		Type:       aType,
		TemplateId: templateId,
		Params:     params,
	}
}
