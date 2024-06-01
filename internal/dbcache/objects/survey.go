package objects

type Survey struct {
	Id int
}

func NewSurvey(id int) *Survey {
	return &Survey{
		Id: id,
	}
}
