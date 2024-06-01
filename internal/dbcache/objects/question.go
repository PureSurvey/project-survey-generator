package objects

import (
	"project-survey-generator/internal/enums"
)

type Question struct {
	Id             int
	Type           enums.QuestionType
	SurveyId       int
	OrderNumber    int
	QuestionLineId int
}

func NewQuestion(id int, qType enums.QuestionType, surveyId int,
	orderNumber int, questionLineId int) *Question {
	return &Question{
		Id:             id,
		Type:           qType,
		SurveyId:       surveyId,
		OrderNumber:    orderNumber,
		QuestionLineId: questionLineId,
	}
}
