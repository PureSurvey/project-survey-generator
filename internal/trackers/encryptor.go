package trackers

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"project-survey-generator/internal/configuration"
	"project-survey-generator/internal/crypto"
	"project-survey-generator/internal/pools"
	"project-survey-generator/internal/trackers/model"
	"strconv"
	"strings"
)

type Encryptor struct {
	encryptionKey []byte

	sbPool *pools.StringBuilder
}

func NewEncryptor(appConfig *configuration.AppConfiguration, sbPool *pools.StringBuilder) *Encryptor {
	key, _ := hex.DecodeString(appConfig.EncryptionSecret)

	return &Encryptor{encryptionKey: key, sbPool: sbPool}
}

func (e *Encryptor) EncryptTracker(event *model.Tracker) string {
	stringToEncrypt := e.getStringFromEvent(event)

	encryptedBytes, _ := crypto.EncryptWithGCM(stringToEncrypt, e.encryptionKey)
	encrypted := base64.URLEncoding.EncodeToString(encryptedBytes)

	return encrypted
}

func (e *Encryptor) getStringFromEvent(event *model.Tracker) string {
	params := [8]string{}

	params[0] = strconv.Itoa(int(event.EventType))
	params[1] = strconv.Itoa(event.UnitId)
	params[2] = strconv.FormatInt(event.ValidTo, 10)

	if event.ValidSurveys != nil && len(event.ValidSurveys) > 0 {
		validSurveysString := join(event.ValidSurveys, ";")
		params[3] = validSurveysString
	}

	if event.ValidQuestions != nil && len(event.ValidQuestions) > 0 {
		validQuestionsString := join(event.ValidQuestions, ";")
		params[4] = validQuestionsString
	}

	if event.ValidQuestionsWithAnswers != nil && len(event.ValidQuestionsWithAnswers) > 0 {
		var validQuestionsWithAnswers []string
		for key, val := range event.ValidQuestionsWithAnswers {
			questionWithAnswers := fmt.Sprintf("%v:%v", key, join(val, ";"))
			validQuestionsWithAnswers = append(validQuestionsWithAnswers, questionWithAnswers)
		}

		params[5] = strings.Join(validQuestionsWithAnswers, "/")
	}

	return strings.Join(params[:], ",")
}

func join(arr []int, sep string) string {
	var stringArr []string
	for _, el := range arr {
		stringArr = append(stringArr, strconv.Itoa(el))
	}

	return strings.Join(stringArr, sep)
}
