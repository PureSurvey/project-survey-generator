package objects

type Template struct {
	Id            int
	Code          *TemplateCode
	DefaultParams map[string]string
}

type TemplateCode struct {
	Unit           string
	Scripts        string
	Styles         string
	Question       string
	QuestionOption string
	NextButton     string
	PrevButton     string
	Completed      string

	NextQuestionClass string
	StartSurveyClass  string
	EndSurveyClass    string
	EndUnitClass      string
}

func NewTemplate(id int, code *TemplateCode, defaultParams map[string]string) *Template {
	return &Template{
		Id:            id,
		Code:          code,
		DefaultParams: defaultParams,
	}
}
