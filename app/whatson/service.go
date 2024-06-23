package whatson

import "gitlab.com/vextasy/claude/whatson/domain"

func NewServices(suggestionSvc domain.SuggestionSvc, tvDbSvc domain.TvDbSvc) domain.SuggestionServices {
	return domain.SuggestionServices{
		Suggestion: suggestionSvc,
	}
}
