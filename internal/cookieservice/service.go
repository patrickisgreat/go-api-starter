package cookieservice

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/patrickisgreat/go-api-starter/internal/cookierepo"
)

type CookieService struct {
	repo *cookierepo.CookieRepository
}

func NewCookieService(repo *cookierepo.CookieRepository) *CookieService {
	return &CookieService{
		repo: repo,
	}
}

func (cs *CookieService) GetCookie(ctx context.Context) (string, twirp.Error) {
	cookie, err := cs.repo.Gimme(ctx)
	if err != nil {
		return "", twirp.InternalError(err.Error())
	}

	return cookie, nil
}
