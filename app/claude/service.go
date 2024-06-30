package claude

import "gitlab.com/vextasy/claude/whatson/domain"

func NewServices(verbose bool, tvDbSvc domain.TvDbSvc) domain.SuggestionServices {
	return domain.SuggestionServices{
		Suggestion: NewSuggestionSvc(verbose, tvDbSvc),
	}
}
