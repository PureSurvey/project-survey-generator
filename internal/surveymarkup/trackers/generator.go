package trackers

import (
	"project-survey-generator/internal/configuration"
	"project-survey-generator/internal/enums"
	"project-survey-generator/internal/pools"
	"project-survey-generator/internal/surveymarkup/model"
	"strconv"
	"strings"
	"time"
)

const (
	TrackerTtl = 30 * time.Minute
)

type Generator struct {
	appConfig *configuration.AppConfiguration
	sbPool    *pools.StringBuilderPool
	encryptor *Encryptor
}

func NewGenerator(appConfig *configuration.AppConfiguration, encryptor *Encryptor, sbPool *pools.StringBuilderPool) *Generator {
	return &Generator{appConfig: appConfig, encryptor: encryptor, sbPool: sbPool}
}

func (tg *Generator) GenerateViewTracker(unitId int) string {
	validTo := time.Now().UTC().Add(TrackerTtl)

	event := &model.Event{
		EventType: enums.ETUnitView,
		UnitId:    unitId,
		ValidTo:   validTo.Unix(),
	}

	return tg.generateTracker(event)
}

func (tg *Generator) GenerateQuestionViewTracker(unitId int, questionIds []int) string {
	validTo := time.Now().UTC().Add(TrackerTtl)

	event := &model.Event{
		EventType:      enums.ETQuestionView,
		UnitId:         unitId,
		ValidTo:        validTo.Unix(),
		ValidQuestions: questionIds,
	}

	return tg.generateTracker(event)
}

func (tg *Generator) GenerateQuestionAnswerTracker(unitId int, optionIdsByQuestionIds map[int][]int) string {
	validTo := time.Now().UTC().Add(TrackerTtl)

	event := &model.Event{
		EventType:                 enums.ETQuestionAnswer,
		UnitId:                    unitId,
		ValidTo:                   validTo.Unix(),
		ValidQuestionsWithAnswers: optionIdsByQuestionIds,
	}

	return tg.generateTracker(event)
}

func (tg *Generator) GenerateSurveyStartTracker(unitId int, surveyIds []int) string {
	validTo := time.Now().UTC().Add(TrackerTtl)

	event := &model.Event{
		EventType:    enums.ETSurveyStart,
		UnitId:       unitId,
		ValidTo:      validTo.Unix(),
		ValidSurveys: surveyIds,
	}

	return tg.generateTracker(event)
}

func (tg *Generator) GenerateSurveyEndTracker(unitId int, surveyIds []int) string {
	validTo := time.Now().UTC().Add(TrackerTtl)

	event := &model.Event{
		EventType:    enums.ETSurveyEnd,
		UnitId:       unitId,
		ValidTo:      validTo.Unix(),
		ValidSurveys: surveyIds,
	}

	return tg.generateTracker(event)
}

func (tg *Generator) GenerateUnitEndTracker(unitId int) string {
	validTo := time.Now().UTC().Add(TrackerTtl)

	event := &model.Event{
		EventType: enums.ETUnitEnd,
		UnitId:    unitId,
		ValidTo:   validTo.Unix(),
	}

	return tg.generateTracker(event)
}

func (tg *Generator) generateTracker(event *model.Event) string {
	sb := tg.sbPool.Get()
	defer tg.sbPool.Put(sb)

	tg.writeTrackerStart(sb, event.EventType)
	tg.writeQueryParam(sb, PUnitId, strconv.Itoa(event.UnitId))
	tg.writeQueryParam(sb, PValidTo, strconv.FormatInt(event.ValidTo, 10))

	encrypted := tg.encryptor.EncryptEvent(event)
	tg.writeQueryParam(sb, PEncryptedEvent, encrypted)

	return sb.String()
}

func (tg *Generator) writeTrackerStart(sb *strings.Builder, eventType enums.EventType) {
	sb.WriteString("/")
	sb.WriteString(tg.appConfig.EventsEndpoint)
	sb.WriteString("?")
	sb.WriteString(PEventType)
	sb.WriteString("=")
	sb.WriteString(strconv.Itoa(int(eventType)))
}

func (tg *Generator) writeQueryParam(sb *strings.Builder, name string, value string) {
	sb.WriteString("&")
	sb.WriteString(name)
	sb.WriteString("=")
	sb.WriteString(value)
}
