package surveymarkup

import (
	"errors"
	"html"
	"math/rand"
	"project-survey-generator/internal/configuration"
	"project-survey-generator/internal/dbcache"
	"project-survey-generator/internal/dbcache/objects"
	"project-survey-generator/internal/enums"
	compilederrors "project-survey-generator/internal/errors"
	"project-survey-generator/internal/localisation"
	"project-survey-generator/internal/pools"
	appearance2 "project-survey-generator/internal/surveymarkup/appearance"
	"project-survey-generator/internal/surveymarkup/macros"
	"project-survey-generator/internal/surveymarkup/minifier"
	"project-survey-generator/internal/surveymarkup/trackers"
	"sort"
	"strconv"
	"strings"
)

type Generator struct {
	dbRepo            *dbcache.Repo
	trackersGenerator *trackers.Generator

	stringBuilderPool *pools.StringBuilder

	appConfiguration *configuration.AppConfiguration
	minifier         *minifier.Service
}

const optionsNameLength = 8
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewGenerator(dbRepo *dbcache.Repo, trackersGenerator *trackers.Generator, minifier *minifier.Service, stringBuilderPool *pools.StringBuilder, appConfig *configuration.AppConfiguration) *Generator {
	return &Generator{dbRepo: dbRepo, trackersGenerator: trackersGenerator, minifier: minifier, stringBuilderPool: stringBuilderPool, appConfiguration: appConfig}
}

func (g *Generator) Generate(unit *objects.Unit, surveys []*objects.Survey, language string) (string, error) {
	appearance := g.dbRepo.GetAppearanceById(unit.AppearanceId)
	if appearance == nil {
		return "", errors.New("")
	}

	template := g.dbRepo.GetTemplateById(appearance.TemplateId)
	if template == nil {
		return "", errors.New("")
	}

	templateCode := template.Code.Unit

	sb := g.stringBuilderPool.Get()
	defer g.stringBuilderPool.Put(sb)

	var surveyIds, questionIds []int
	optionIdsByQuestionIds := map[int][]int{}

	for i, survey := range surveys {
		surveyIds = append(surveyIds, survey.Id)

		questions := g.dbRepo.GetQuestionsBySurveyId(survey.Id)
		sort.Slice(questions, func(prev, cur int) bool {
			return questions[prev].OrderNumber < questions[cur].OrderNumber
		})

		var surveyEndCode string
		if i != 0 {
			surveyEndCode = g.generateSurveyEnd(template, questions[0], language)
		}
		sb.WriteString(surveyEndCode)

		for j, question := range questions {
			options := g.dbRepo.GetOptionsByQuestionId(question.Id)
			if options == nil {
				return "", compilederrors.NoOptionsInQuestion
			}

			questionIds = append(questionIds, question.Id)
			optionIdsByQuestionIds[question.Id] = Map(options, func(t *objects.Option) int {
				return t.Id
			})

			var questionCode string
			var err error
			if j == len(questions)-1 {
				questionCode, err = g.generateQuestion(question, options, nil, template, false, language)
			} else {
				questionCode, err = g.generateQuestion(question, options, questions[j+1], template, i == 0, language)
			}

			if err != nil {
				return "", err
			}
			sb.WriteString(questionCode)
		}
	}

	surveyEndCode := g.generateSurveyEnd(template, nil, language)
	sb.WriteString(surveyEndCode)

	unitEndCode := g.generateUnitEnd(template, language)
	sb.WriteString(unitEndCode)

	styles := template.Code.Styles
	for key, tmParamVal := range template.DefaultParams {
		if appParamVal, ok := appearance.Params[key]; ok {
			styles = strings.Replace(styles, key, appParamVal, -1)
		} else {
			styles = strings.Replace(styles, key, tmParamVal, -1)
		}
	}

	viewTracker := g.trackersGenerator.GenerateViewTracker(unit.Id)
	questionViewTracker := g.trackersGenerator.GenerateQuestionViewTracker(unit.Id, questionIds)
	questionAnswerTracker := g.trackersGenerator.GenerateQuestionAnswerTracker(unit.Id, optionIdsByQuestionIds)
	surveyStartTracker := g.trackersGenerator.GenerateSurveyStartTracker(unit.Id, surveyIds)
	surveyEndTracker := g.trackersGenerator.GenerateSurveyEndTracker(unit.Id, surveyIds)
	unitEndTracker := g.trackersGenerator.GenerateUnitEndTracker(unit.Id)

	scriptCode := template.Code.Scripts
	scriptCode = strings.Replace(scriptCode, macros.QuestionViewTracker, questionViewTracker, -1)
	scriptCode = strings.Replace(scriptCode, macros.QuestionAnswerTracker, questionAnswerTracker, 1)
	scriptCode = strings.Replace(scriptCode, macros.SurveyStartTracker, surveyStartTracker, 1)
	scriptCode = strings.Replace(scriptCode, macros.SurveyEndTracker, surveyEndTracker, 1)
	scriptCode = strings.Replace(scriptCode, macros.UnitEndTracker, unitEndTracker, 1)
	scriptCode = strings.Replace(scriptCode, macros.ProceederUrl, g.appConfiguration.SurveyProceederUrl, 1)

	if unit.HideAfterNoSurveys {
		scriptCode = strings.Replace(scriptCode, macros.HideUnitAfterEnd, "1", -1)
	} else {
		scriptCode = strings.Replace(scriptCode, macros.HideUnitAfterEnd, "0", -1)
	}

	templateCode = strings.Replace(templateCode, macros.SurveyHeaderText, localisation.HeaderTextByLanguage[language], 1)
	templateCode = strings.Replace(templateCode, macros.ViewPixelUrl, g.appConfiguration.SurveyProceederUrl+viewTracker, 1)
	templateCode = strings.Replace(templateCode, macros.Styles, styles, 1)
	templateCode = strings.Replace(templateCode, macros.UnitItems, sb.String(), 1)

	switch appearance.Type {
	case enums.ATLeft:
		templateCode = strings.Replace(appearance2.Left, macros.SurveyPlacement, templateCode, 1)
		scriptCode += appearance2.LeftScript
	case enums.ATRight:
		templateCode = strings.Replace(appearance2.Right, macros.SurveyPlacement, templateCode, 1)
		scriptCode += appearance2.RightScript
	case enums.ATOverlay:
		templateCode = strings.Replace(appearance2.Overlay, macros.SurveyPlacement, templateCode, 1)
		scriptCode += appearance2.OverlayScript
	default:
		break
	}

	scriptCode = g.minifier.MinifyJs(scriptCode)
	scriptCode = strings.Replace(scriptCode, "\"", "\\\"", -1)
	scriptCode = strings.Replace(scriptCode, "\n", " ", -1)
	scriptCode = strings.Replace(scriptCode, "\r", " ", -1)
	scriptCode = html.EscapeString(scriptCode)

	templateCode = g.minifier.Minify(templateCode)
	templateCode = strings.Replace(templateCode, "\"", "\\\"", -1)
	templateCode = strings.Replace(templateCode, "\n", " ", -1)
	templateCode = strings.Replace(templateCode, "\r", " ", -1)
	templateCode = html.EscapeString(templateCode)

	templateCode = strings.Replace(appearance2.ResponseScript, macros.Html, templateCode, 1)
	templateCode = strings.Replace(templateCode, macros.UnitId, strconv.Itoa(unit.Id), 1)
	templateCode = strings.Replace(templateCode, macros.Script, scriptCode, 1)

	return templateCode, nil
}

func (g *Generator) generateQuestion(question *objects.Question, options []*objects.Option, nextQuestion *objects.Question, template *objects.Template, isFirstSurvey bool, language string) (string, error) {
	questionCode := template.Code.Question

	questionLineTranslations := g.dbRepo.GetTranslationsByQuestionLineId(question.QuestionLineId)
	if questionLineTranslations == nil {
		return "", compilederrors.TranslationNotFound
	}

	questionLineTranslation := questionLineTranslations[language]
	if questionLineTranslation == nil {
		return "", compilederrors.TranslationNotFound
	}

	questionCode = strings.Replace(questionCode, macros.QuestionBlockText, html.EscapeString(questionLineTranslation.Translation), 1)

	// next button
	nextButtonCode := template.Code.NextButton
	nextButtonText := localisation.NextButtonTextByLanguage[language]
	var nextButtonClasses []string

	if question.OrderNumber == 1 {
		nextButtonClasses = append(nextButtonClasses, template.Code.StartSurveyClass)
	}

	if nextQuestion == nil {
		nextButtonClasses = append(nextButtonClasses, template.Code.EndSurveyClass)
		nextButtonCode = strings.Replace(nextButtonCode, macros.NextQuestionId, "", 1)
	} else {
		nextButtonClasses = append(nextButtonClasses, template.Code.NextQuestionClass)
		nextButtonCode = strings.Replace(nextButtonCode, macros.NextQuestionId, strconv.Itoa(nextQuestion.Id), 1)
	}

	nextButtonClassesString := strings.Join(nextButtonClasses, " ")
	nextButtonCode = strings.Replace(nextButtonCode, macros.NextButtonClasses, nextButtonClassesString, 1)
	nextButtonCode = strings.Replace(nextButtonCode, macros.NextButtonText, nextButtonText, 1)
	nextButtonCode = strings.Replace(nextButtonCode, macros.SurveyId, strconv.Itoa(question.SurveyId), 1)
	nextButtonCode = strings.Replace(nextButtonCode, macros.QuestionId, strconv.Itoa(question.Id), 1)

	var prevButtonCode string
	if !isFirstSurvey && question.OrderNumber != 1 {
		prevButtonCode = template.Code.PrevButton
		prevButtonText := localisation.PreviousButtonTextByLanguage[language]
		prevButtonCode = strings.Replace(prevButtonCode, macros.PrevButtonText, prevButtonText, 1)
	}

	questionCode = strings.Replace(questionCode, macros.NextButton, nextButtonCode, 1)
	questionCode = strings.Replace(questionCode, macros.PrevButton, prevButtonCode, 1)

	// question options
	optionsType := ""
	if question.Type == enums.QTRadiobutton {
		optionsType = "radio"
	} else if question.Type == enums.QTCheckbox {
		optionsType = "checkbox"
	} else {
		return "", compilederrors.UnknownQuestionType
	}

	sort.Slice(options, func(prev, cur int) bool {
		return options[prev].OrderNumber < options[cur].OrderNumber
	})

	sb := g.stringBuilderPool.Get()
	defer g.stringBuilderPool.Put(sb)

	optionsName := g.generateOptionsName()
	for _, option := range options {
		optionCode := template.Code.QuestionOption

		optionTranslations := g.dbRepo.GetTranslationsByOptionId(option.Id)
		if optionTranslations == nil {
			return "", compilederrors.TranslationNotFound
		}

		optionTranslation := optionTranslations[language]
		if optionTranslation == nil {
			return "", compilederrors.TranslationNotFound
		}

		optionCode = strings.Replace(optionCode, macros.OptionsType, optionsType, 1)
		optionCode = strings.Replace(optionCode, macros.OptionsValue, strconv.Itoa(option.Id), 1)
		optionCode = strings.Replace(optionCode, macros.OptionsName, optionsName, 1)
		optionCode = strings.Replace(optionCode, macros.OptionsText, html.EscapeString(optionTranslation.Translation), 1)

		sb.WriteString(optionCode)
	}

	questionCode = strings.Replace(questionCode, macros.QuestionOptions, sb.String(), 1)

	return questionCode, nil
}

func (g *Generator) generateSurveyEnd(template *objects.Template, nextQuestion *objects.Question, language string) string {
	surveyEndCode := template.Code.Completed
	surveyEndCode = strings.Replace(surveyEndCode, macros.CompletedBlockText, localisation.SurveyEndTextByLanguage[language], 1)

	// next button
	nextButtonCode := template.Code.NextButton
	nextButtonText := localisation.NextSurveyButtonTextByLanguage[language]
	var nextButtonClasses []string
	if nextQuestion == nil {
		nextButtonClasses = append(nextButtonClasses, template.Code.EndUnitClass)
		nextButtonCode = strings.Replace(nextButtonCode, macros.NextQuestionId, "", 1)

		nextButtonCode = strings.Replace(nextButtonCode, macros.SurveyId, "", 1)
		nextButtonCode = strings.Replace(nextButtonCode, macros.QuestionId, "", 1)
	} else {
		nextButtonCode = strings.Replace(nextButtonCode, macros.SurveyId, strconv.Itoa(nextQuestion.SurveyId), 1)
		nextButtonCode = strings.Replace(nextButtonCode, macros.QuestionId, strconv.Itoa(nextQuestion.Id), 1)
	}

	nextButtonClassesString := strings.Join(nextButtonClasses, " ")
	nextButtonCode = strings.Replace(nextButtonCode, macros.NextButtonClasses, nextButtonClassesString, 1)
	nextButtonCode = strings.Replace(nextButtonCode, macros.NextButtonText, nextButtonText, 1)

	// prev button
	prevButtonCode := template.Code.PrevButton
	prevButtonText := localisation.PreviousButtonTextByLanguage[language]
	prevButtonCode = strings.Replace(prevButtonCode, macros.PrevButtonText, prevButtonText, 1)

	surveyEndCode = strings.Replace(surveyEndCode, macros.NextButton, nextButtonCode, 1)
	surveyEndCode = strings.Replace(surveyEndCode, macros.PrevButton, prevButtonCode, 1)

	return surveyEndCode
}

func (g *Generator) generateUnitEnd(template *objects.Template, language string) string {
	unitEndCode := template.Code.Completed
	unitEndCode = strings.Replace(unitEndCode, macros.CompletedBlockText, localisation.UnitEndTextByLanguage[language], 1)

	unitEndCode = strings.Replace(unitEndCode, macros.NextButton, "", 1)
	unitEndCode = strings.Replace(unitEndCode, macros.PrevButton, "", 1)

	return unitEndCode
}

func (g *Generator) generateOptionsName() string {
	b := make([]byte, optionsNameLength)
	b[0] = 'p'
	b[1] = '-'
	for i := 2; i < optionsNameLength; i++ {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
