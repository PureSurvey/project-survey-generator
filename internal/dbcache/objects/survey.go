package objects

type Survey struct {
	Id   int
	Name string
}

func NewSurvey(id int, name string) *Survey {
	return &Survey{
		Id:   id,
		Name: name,
	}
}
