package dbcache

import (
	"project-survey-generator/internal/dbcache/contracts"
	"project-survey-generator/internal/dbcache/objects"
	"slices"
	"time"
)

const (
	sleepTime           = 30 * time.Second
	storedProcedureName = "generator_01"
)

type Repo struct {
	reader contracts.IReader
	cache  *Cache
}

func NewRepo(reader contracts.IReader) *Repo {
	return &Repo{reader: reader}
}

func (r *Repo) RunReloadCycle() {
	for {
		r.reload()
		time.Sleep(sleepTime)
	}
}

func (r *Repo) reload() {
	err := r.reader.Connect()
	defer r.reader.CloseConnection()
	if err != nil {
		return
	}

	newCache := &Cache{
		Units:       map[int]*objects.Unit{},
		Surveys:     map[int]*objects.Survey{},
		Templates:   map[int]*objects.Template{},
		Appearances: map[int]*objects.Appearance{},
		Questions:   map[int]*objects.Question{},
		Options:     map[int]*objects.Option{},

		SurveysByUnit:     map[int][]*objects.Survey{},
		QuestionsBySurvey: map[int][]*objects.Question{},
		OptionsByQuestion: map[int][]*objects.Option{},

		TranslationsByQuestionLine: map[int]map[string]*objects.Translation{},
		TranslationsByOption:       map[int]map[string]*objects.Translation{},
	}

	res, err := r.reader.GetStoredProcedureResult(storedProcedureName)
	i := 0
	for cont := true; cont; cont = res.NextResultSet() {
		switch i {
		case 0:
			err := newCache.fillUnits(res)
			if err != nil {
				return
			}
			break
		case 1:
			err := newCache.fillSurveys(res)
			if err != nil {
				return
			}
			break
		case 2:
			err := newCache.fillSurveyInUnits(res)
			if err != nil {
				return
			}
			break
		case 3:
			err := newCache.fillTemplates(res)
			if err != nil {
				return
			}
			break
		case 4:
			err := newCache.fillAppearances(res)
			if err != nil {
				return
			}
			break
		case 5:
			err := newCache.fillQuestions(res)
			if err != nil {
				return
			}
			break
		case 6:
			err := newCache.fillOptions(res)
			if err != nil {
				return
			}
			break
		case 7:
			err := newCache.fillQuestionTranslations(res)
			if err != nil {
				return
			}
			break
		case 8:
			err := newCache.fillOptionTranslations(res)
			if err != nil {
				return
			}
			break
		}
		i++
	}

	r.cache = newCache
}

func (r *Repo) GetUnitById(id int) *objects.Unit {
	return r.cache.Units[id]
}

func (r *Repo) GetSurveysByUnitId(id int) []*objects.Survey {
	return r.cache.SurveysByUnit[id]
}

func (r *Repo) GetQuestionsBySurveyId(id int) []*objects.Question {
	return r.cache.QuestionsBySurvey[id]
}

func (r *Repo) GetOptionsByQuestionId(id int) []*objects.Option {
	return r.cache.OptionsByQuestion[id]
}

func (r *Repo) GetUnitSurveysWithIds(unitId int, surveysIds []int) []*objects.Survey {
	surveys := r.cache.SurveysByUnit[unitId]
	var matchedSurveys []*objects.Survey
	for _, survey := range surveys {
		if slices.Contains(surveysIds, survey.Id) {
			matchedSurveys = append(matchedSurveys, survey)
		}
	}

	return matchedSurveys
}

func (r *Repo) GetAppearanceById(id int) *objects.Appearance {
	return r.cache.Appearances[id]
}

func (r *Repo) GetTemplateById(id int) *objects.Template {
	return r.cache.Templates[id]
}

func (r *Repo) GetTranslationsByQuestionLineId(id int) map[string]*objects.Translation {
	return r.cache.TranslationsByQuestionLine[id]
}

func (r *Repo) GetTranslationsByOptionId(id int) map[string]*objects.Translation {
	return r.cache.TranslationsByOption[id]
}
