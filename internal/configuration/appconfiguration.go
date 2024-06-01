package configuration

type AppConfiguration struct {
	DbConnectionString  string `json:"dbConnectionString"`
	SurveyGeneratorPort int    `json:"surveyGeneratorPort"`
	SurveyProceederUrl  string `json:"surveyProceederUrl"`
	EventsEndpoint      string `json:"eventsEndpoint"`
	EncryptionSecret    string `json:"encryptionSecret"` // 32 bytes
}
