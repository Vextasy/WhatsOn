package domain

import "context"

// Suggestion Services - Services used by the whatson command.
type SuggestionServices struct {
	Suggestion SuggestionSvc
}

type SuggestionSvc interface {
	GetSuggestions(ctx context.Context, desire string) ([]Suggestion, error)
}

// Services used by the TvDB Service.
type TvDbServices struct {
	TvDb TvDbSvc
}

type TvDbSvc interface {
	GetTvProgrammesXml(ctx context.Context, dateFrom string, dateTo string) (string, error)
}
