package claude

import "gitlab.com/vextasy/claude/whatson/domain"

func NewServices(tvDbSvc domain.TvDbSvc) domain.SuggestionServices {
	return domain.SuggestionServices{
		Suggestion: NewSuggestionSvc(tvDbSvc),
	}
}
