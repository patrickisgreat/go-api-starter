package api

import (
	"context"
	"math/rand"

	"github.com/twitchtv/twirp"

	"github.com/patrickisgreat/pb-go-api-starter/internal/appserver"
	proto "github.com/patrickisgreat/pb-go-api-starter/internal/generated/pb_go_api_starter/api"
)

func GetCookieFlakyHandler(ctx context.Context, serverContext *appserver.ServerContext, request *proto.GetCookieFlakyRequest) (*proto.GetCookieFlakyResponse, error) {
	// validate user session
	err := validateUserSession(request.UserSession)
	if err != nil {
		return nil, err
	}

	if request.GetFail() || rand.Intn(2) == 0 {
		return nil, twirp.InternalError("flaky error")
	}

	cookieService := serverContext.CookieService

	cookie, err := cookieService.GetCookie(ctx)

	return &proto.GetCookieFlakyResponse{
		FortuneCookieMessage: cookie,
	}, err
}
