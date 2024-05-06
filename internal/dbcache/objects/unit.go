package objects

type Unit struct {
	Id           int
	AppearanceId int

	HideAfterNoSurveys    bool
	MessageAfterNoSurveys string
}

func NewUnit(id int, appearanceId int,
	hideAfterNoSurveys bool, message string) *Unit {
	return &Unit{
		Id:                    id,
		AppearanceId:          appearanceId,
		HideAfterNoSurveys:    hideAfterNoSurveys,
		MessageAfterNoSurveys: message,
	}
}
