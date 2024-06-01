package configuration

type AppConfiguration struct {
	SurveyGeneratorPort  int                   `json:"surveyGeneratorPort"`
	SurveyProceederUrl   string                `json:"surveyProceederUrl"`
	EventsEndpoint       string                `json:"eventsEndpoint"`
	EncryptionSecret     string                `json:"encryptionSecret"` // 32 bytes
	DbCacheConfiguration *DbCacheConfiguration `json:"dbCacheConfiguration"`
}

type DbCacheConfiguration struct {
	ConnectionRetryCount     int    `json:"connectionRetryCount"`
	ConnectionRetrySleepTime int    `json:"connectionRetryTimeout"`
	ConnectionString         string `json:"connectionString"`
	StoredProcedure          string `json:"storedProcedure"`
	ReloadSleepTime          int    `json:"reloadSleepTime"`
}
