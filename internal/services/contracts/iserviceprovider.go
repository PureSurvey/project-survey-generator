package contracts

import (
	"project-survey-generator/internal/dbcache"
	"project-survey-generator/internal/surveymarkup"
	"project-survey-generator/internal/surveymarkup/minifier"
)

type IServiceProvider interface {
	GetDbRepo() *dbcache.Repo
	GetMinifier() *minifier.Service
	GetGenerator() *surveymarkup.Generator
}
