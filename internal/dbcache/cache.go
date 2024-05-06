package dbcache

import (
	"database/sql"
	"project-survey-generator/internal/dbcache/objects"
	"project-survey-generator/internal/enums"
)

type Cache struct {
	Units       map[int]*objects.Unit
	Surveys     map[int]*objects.Survey
	Templates   map[int]*objects.Template
	Appearances map[int]*objects.Appearance
	Questions   map[int]*objects.Question
	Options     map[int]*objects.Option

	SurveysByUnit map[int][]*objects.Survey

	TranslationsByQuestionLine map[int]map[string]*objects.Translation
	TranslationsByOption       map[int]map[string]*objects.Translation
}

func (c *Cache) fillUnits(rows *sql.Rows) error {
	for rows.Next() {
		var id, appearanceId int
		var hideAfterNoSurveys bool
		var message string

		err := rows.Scan(&id, &appearanceId, &hideAfterNoSurveys, &message)
		if err != nil {
			return err
		}

		unit := objects.NewUnit(id, appearanceId, hideAfterNoSurveys, message)
		c.Units[id] = unit
	}

	return nil
}

func (c *Cache) fillSurveys(rows *sql.Rows) error {
	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			return err
		}

		survey := objects.NewSurvey(id, name)
		c.Surveys[id] = survey
	}

	return nil
}

func (c *Cache) fillSurveyInUnits(rows *sql.Rows) error {
	for rows.Next() {
		var surveyId, surveyUnitId int

		err := rows.Scan(&surveyId, &surveyUnitId)
		if err != nil {
			return err
		}

		survey := c.Surveys[surveyId]
		if survey != nil {
			c.SurveysByUnit[surveyUnitId] = append(c.SurveysByUnit[surveyUnitId], survey)
		}
	}

	return nil
}

func (c *Cache) fillTemplates(rows *sql.Rows) error {
	for rows.Next() {
		var id int
		var code, defaultParams string

		err := rows.Scan(&id, &code, &defaultParams)
		if err != nil {
			return err
		}

		template := objects.NewTemplate(id, code, defaultParams)
		c.Templates[id] = template
	}

	return nil
}

func (c *Cache) fillAppearances(rows *sql.Rows) error {
	for rows.Next() {
		var id, aType, templateId int
		var params string

		err := rows.Scan(&id, &aType, &templateId, &params)
		if err != nil {
			return err
		}

		appearance := objects.NewAppearance(id, enums.EnumAppearanceType(aType), templateId, params)
		c.Appearances[id] = appearance
	}

	return nil
}

func (c *Cache) fillQuestions(rows *sql.Rows) error {
	for rows.Next() {
		var id, qType, surveyId, orderNumber, questionLineId int

		err := rows.Scan(&id, &qType, &surveyId, &orderNumber, &questionLineId)
		if err != nil {
			return err
		}

		question := objects.NewQuestion(id, enums.EnumQuestionType(qType), surveyId, orderNumber, questionLineId)
		c.Questions[id] = question
	}

	return nil
}

func (c *Cache) fillOptions(rows *sql.Rows) error {
	for rows.Next() {
		var id, questionId, orderNumber int

		err := rows.Scan(&id, &questionId, &orderNumber)
		if err != nil {
			return err
		}

		option := objects.NewOption(id, questionId, orderNumber)
		c.Options[id] = option
	}

	return nil
}

func (c *Cache) fillQuestionTranslations(rows *sql.Rows) error {
	for rows.Next() {
		var id, questionLineId int
		var lang, translationLine string

		err := rows.Scan(&id, &lang, &translationLine, &questionLineId)
		if err != nil {
			return err
		}

		translation := objects.NewTranslation(id, lang, translationLine, questionLineId)

		if c.TranslationsByQuestionLine[questionLineId] == nil {
			c.TranslationsByQuestionLine[questionLineId] = map[string]*objects.Translation{}
		}

		c.TranslationsByQuestionLine[questionLineId][lang] = translation
	}

	return nil
}

func (c *Cache) fillOptionTranslations(rows *sql.Rows) error {
	for rows.Next() {
		var id, optionId int
		var lang, translationLine string

		err := rows.Scan(&id, &lang, &translationLine, &optionId)
		if err != nil {
			return err
		}

		translation := objects.NewTranslation(id, lang, translationLine, optionId)

		if c.TranslationsByOption[optionId] == nil {
			c.TranslationsByOption[optionId] = map[string]*objects.Translation{}
		}

		c.TranslationsByOption[optionId][lang] = translation
	}

	return nil
}
