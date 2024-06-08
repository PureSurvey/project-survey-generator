package model

import "project-survey-generator/internal/enums"

type Tracker struct {
	EventType                 enums.EventType
	UnitId                    int
	ValidTo                   int64
	ValidSurveys              []int
	ValidQuestions            []int
	ValidQuestionsWithAnswers map[int][]int
}
