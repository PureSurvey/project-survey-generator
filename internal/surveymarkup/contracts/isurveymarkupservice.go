package contracts

type ISurveyMarkupServer interface {
	GetMarkup(unitId int, surveyIds []int, language string) (string, error)
}
