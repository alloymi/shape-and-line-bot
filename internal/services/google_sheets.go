package services

import (
	"context"
	"encoding/json"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"SnLbot/internal/config"
)

func SaveToGoogleSheet(fullname, email, course string) error {
	cfg := config.Load()

	// Распарсим JSON ключ из env
	var creds map[string]interface{}
	json.Unmarshal([]byte(cfg.GoogleServiceAccountJSON), &creds)

	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(cfg.GoogleServiceAccountJSON)))
	if err != nil {
		return err
	}

	row := []interface{}{fullname, email, course}

	_, err = srv.Spreadsheets.Values.Append(
		cfg.GoogleSpreadsheetID,
		"A1",
		&sheets.ValueRange{Values: [][]interface{}{row}},
	).ValueInputOption("RAW").Do()

	return err
}
