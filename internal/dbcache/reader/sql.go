package reader

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"time"
)

const (
	retryCount     = 3
	retrySleepTime = 10 * time.Second
)

type SqlReader struct {
	connectionString string
	db               *sql.DB
}

func NewSqlReader(connectionString string) *SqlReader {
	return &SqlReader{connectionString: connectionString}
}

func (s *SqlReader) Connect() error {
	var err error
	s.db, err = sql.Open("sqlserver", s.connectionString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
		return err
	}

	ctx := context.Background()
	for i := 0; i < retryCount; i++ {
		err = s.db.PingContext(ctx)
		if err != nil {
			log.Println("Error when connecting to DB:", err.Error(), "Try", i+1, "of", retryCount)
			time.Sleep(retrySleepTime)
		}
	}
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func (s *SqlReader) GetStoredProcedureResult(storedProcedure string) (*sql.Rows, error) {
	return s.db.Query(fmt.Sprintf("exec %v", storedProcedure))
}
