package function

import (
	"context"

	sheets "google.golang.org/api/sheets/v4"
)

type SheetsConnection struct {
	ctx            context.Context
	service        *sheets.Service
	spreadsheet_id string
}

func NewSheetsConnection(ctx context.Context, spreadsheet_id string) *SheetsConnection {
	service, err := sheets.NewService(ctx)
	if err != nil {
		panic(err)
	}
	return &SheetsConnection{
		ctx:            ctx,
		service:        service,
		spreadsheet_id: spreadsheet_id,
	}
}
