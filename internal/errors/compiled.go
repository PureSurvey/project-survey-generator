package errors

import "errors"

var TranslationNotFound = errors.New("translation was not found")
var NoOptionsInQuestion = errors.New("no options in question")
var UnknownQuestionType = errors.New("unknown question type")
