package api

import (
	"context"

	"github.com/patrickisgreat/pb-go-api-starter/internal/appserver"
	proto "github.com/patrickisgreat/pb-go-api-starter/internal/generated/pb_go_api_starter/api"
)

func GetQuoteHandler(ctx context.Context, serverContext *appserver.ServerContext, request *proto.GetQuoteRequest) (*proto.GetQuoteResponse, error) {
	// validate user session
	err := validateUserSession(request.GetUserSession())
	if err != nil {
		return nil, err
	}

	quoteService := serverContext.QuoteService

	quote, err := quoteService.GetQuote(ctx)

	return &proto.GetQuoteResponse{
		Quote: quote,
	}, err
}
