package trackers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"project-survey-generator/internal/configuration"
	"project-survey-generator/internal/pools"
	"project-survey-generator/internal/surveymarkup/model"
	"strconv"
	"strings"
)

type Encryptor struct {
	encryptionKey []byte

	sbPool *pools.StringBuilderPool
}

func NewEncryptor(appConfig *configuration.AppConfiguration, sbPool *pools.StringBuilderPool) *Encryptor {
	key, _ := hex.DecodeString(appConfig.EncryptionSecret)

	return &Encryptor{encryptionKey: key, sbPool: sbPool}
}

func (e *Encryptor) EncryptEvent(event *model.Event) string {
	stringToEncrypt := e.getStringFromEvent(event)

	encrypted, _ := encrypt(stringToEncrypt, e.encryptionKey)

	return encrypted
}

func encrypt(stringToEncrypt string, key []byte) (string, error) {
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryptor) getStringFromEvent(event *model.Event) string {
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
		for key, val := range event.ValidQuestionsWithAnswers {
			validQuestionsString := fmt.Sprintf("%v:%v", key, join(val, ";"))
			params[5] = validQuestionsString
		}
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
