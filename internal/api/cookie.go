package api

import (
	"context"

	"github.com/patrickisgreat/pb-go-api-starter/internal/appserver"
	proto "github.com/patrickisgreat/pb-go-api-starter/internal/generated/pb_go_api_starter/api"
)

func GetCookieHandler(ctx context.Context, serverContext *appserver.ServerContext, request *proto.GetCookieRequest) (*proto.GetCookieResponse, error) {
	// validate user session
	err := validateUserSession(request.GetUserSession())
	if err != nil {
		return nil, err
	}

	cookieService := serverContext.CookieService

	cookie, err := cookieService.GetCookie(ctx)

	return &proto.GetCookieResponse{
		FortuneCookieMessage: cookie,
	}, err
}
