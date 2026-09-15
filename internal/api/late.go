package api

import (
	"context"
	"math/rand"
	"time"

	"github.com/patrickisgreat/pb-go-api-starter/internal/appserver"
	proto "github.com/patrickisgreat/pb-go-api-starter/internal/generated/pb_go_api_starter/api"
)

func GetCookieLateHandler(ctx context.Context, serverContext *appserver.ServerContext, request *proto.GetCookieLateRequest) (*proto.GetCookieLateResponse, error) {
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

	cookieService := serverContext.CookieService

	cookie, err := cookieService.GetCookie(ctx)

	return &proto.GetCookieLateResponse{
		FortuneCookieMessage: cookie,
	}, err
}
