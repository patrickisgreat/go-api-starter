package api

import (
	"context"
	"math/rand"
	"time"

	"github.com/patrickisgreat/pb-go-api-starter/internal/appserver"
	proto "github.com/patrickisgreat/pb-go-api-starter/internal/generated/pb_go_api_starter/api"
)

func GetQuoteLateHandler(ctx context.Context, serverContext *appserver.ServerContext, request *proto.GetQuoteLateRequest) (*proto.GetQuoteLateResponse, error) {
	// validate user session
	err := validateUserSession(request.UserSession)
	if err != nil {
		return nil, err
	}

	delay := request.GetDelayMs().GetValue()
	if delay == 0 {
		delay = int32(rand.Float64()*900 + 100)
	}

	time.Sleep(time.Duration(delay) * time.Millisecond)

	quoteService := serverContext.QuoteService

	quote, err := quoteService.GetQuote(ctx)

	return &proto.GetQuoteLateResponse{
		Quote: quote,
	}, err
}
