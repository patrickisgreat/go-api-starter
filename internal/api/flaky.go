package api

import (
	"context"
	"math/rand"

	"github.com/twitchtv/twirp"

	"github.com/patrickisgreat/pb-go-api-starter/internal/appserver"
	proto "github.com/patrickisgreat/pb-go-api-starter/internal/generated/pb_go_api_starter/api"
)

func GetQuoteFlakyHandler(ctx context.Context, serverContext *appserver.ServerContext, request *proto.GetQuoteFlakyRequest) (*proto.GetQuoteFlakyResponse, error) {
	// validate user session
	err := validateUserSession(request.UserSession)
	if err != nil {
		return nil, err
	}

	if request.GetFail() || rand.Intn(2) == 0 {
		return nil, twirp.InternalError("flaky error")
	}

	quoteService := serverContext.QuoteService

	quote, err := quoteService.GetQuote(ctx)

	return &proto.GetQuoteFlakyResponse{
		Quote: quote,
	}, err
}
