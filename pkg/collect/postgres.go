package collect

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	troubleshootv1beta1 "github.com/replicatedhq/troubleshoot/pkg/apis/troubleshoot/v1beta1"
)

type PostgresOutput map[string][]byte

type DatabaseConnection struct {
	IsConnected bool   `json:"isConnected"`
	Error       string `json:"error"`
	Version     string `json:"version"`
}

func Postgres(ctx *Context, postgresCollector *troubleshootv1beta1.Postgres) ([]byte, error) {
	databaseConnection := DatabaseConnection{}

	db, err := sql.Open("postgres", postgresCollector.URI)
	if err != nil {
		databaseConnection.Error = err.Error()
	} else {
		query := `select version()`
		row := db.QueryRow(query)
		version := ""
		if err := row.Scan(&version); err != nil {
			databaseConnection.Error = err.Error()
		} else {
			databaseConnection.IsConnected = true
			databaseConnection.Version = version
		}
	}

	b, err := json.Marshal(databaseConnection)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal database connection")
	}

	collectorName := postgresCollector.CollectorName
	if collectorName == "" {
		collectorName = "postgres"
	}

	postgresOutput := map[string][]byte{
		fmt.Sprintf("postgres/%s.json", collectorName): b,
	}

	bb, err := json.Marshal(postgresOutput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal postgres output")
	}

	return bb, nil
}