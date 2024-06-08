package services

import (
	"project-survey-generator/internal/configuration"
	"project-survey-generator/internal/dbcache"
	"project-survey-generator/internal/dbcache/reader"
	"project-survey-generator/internal/pools"
	"project-survey-generator/internal/services/contracts"
	"project-survey-generator/internal/surveymarkup"
	"project-survey-generator/internal/surveymarkup/minifier"
	trackers2 "project-survey-generator/internal/trackers"
)

type Provider struct {
	dbRepo            *dbcache.Repo
	generator         *surveymarkup.Generator
	trackersGenerator *trackers2.Generator
	encryptor         *trackers2.Encryptor
	minifier          *minifier.Service
}

func NewProvider(appConfiguration *configuration.AppConfiguration) contracts.IServiceProvider {
	dbReader := reader.NewSqlReader(appConfiguration.DbCacheConfiguration)
	dbRepo := dbcache.NewRepo(appConfiguration.DbCacheConfiguration, dbReader)
	go dbRepo.RunReloadCycle()

	sbPool := pools.NewStringBuilderPool()

	encryptor := trackers2.NewEncryptor(appConfiguration, sbPool)
	trackersGenerator := trackers2.NewGenerator(appConfiguration, encryptor, sbPool)
	minifier := minifier.NewService()

	provider := &Provider{
		dbRepo:    dbRepo,
		generator: surveymarkup.NewGenerator(dbRepo, trackersGenerator, minifier, sbPool, appConfiguration),
		minifier:  minifier,
	}

	return provider
}

func (p *Provider) GetDbRepo() *dbcache.Repo {
	return p.dbRepo
}

func (p *Provider) GetGenerator() *surveymarkup.Generator {
	return p.generator
}

func (p *Provider) GetMinifier() *minifier.Service {
	return p.minifier
}
