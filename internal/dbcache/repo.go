package dbcache

import (
	"project-survey-generator/internal/dbcache/contracts"
	"project-survey-generator/internal/dbcache/objects"
)

const (
	storedProcedureName = "generator_01"
)

type Repo struct {
	reader contracts.IReader
	cache  *Cache
}

func NewRepo(reader contracts.IReader) *Repo {
	return &Repo{reader: reader}
}

func (r *Repo) Reload() {
	err := r.reader.Connect()
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

		SurveysByUnit: map[int][]*objects.Survey{},

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
