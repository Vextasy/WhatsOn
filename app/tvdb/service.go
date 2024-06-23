package tvdb

import "gitlab.com/vextasy/claude/whatson/domain"

func NewServices(pathToTvDb string) domain.TvDbServices {
	return domain.TvDbServices{
		TvDb: NewTvDbSvc(pathToTvDb),
	}
}
