package quoteservice

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/patrickisgreat/pb-go-api-starter/internal/quoterepo"
)

type QuoteService struct {
	repo *quoterepo.QuoteRepository
}

func NewQuoteService(repo *quoterepo.QuoteRepository) *QuoteService {
	return &QuoteService{
		repo: repo,
	}
}

func (cs *QuoteService) GetQuote(ctx context.Context) (string, twirp.Error) {
	quote, err := cs.repo.Gimme(ctx)
	if err != nil {
		return "", twirp.InternalError(err.Error())
	}

	return quote, nil
}
