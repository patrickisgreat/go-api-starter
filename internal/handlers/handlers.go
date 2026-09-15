package handlers

import (
	"context"

	"github.com/patrickisgreat/pb-go-api-starter/internal/api"
	"github.com/patrickisgreat/pb-go-api-starter/internal/appserver"
	proto "github.com/patrickisgreat/pb-go-api-starter/internal/generated/pb_go_api_starter/api"
)

type Fortune struct {
	ServerContext *appserver.ServerContext
}

func (h *Fortune) GetCookie(ctx context.Context, request *proto.GetCookieRequest) (*proto.GetCookieResponse, error) {
	return api.GetCookieHandler(ctx, h.ServerContext, request)
}

func (h *Fortune) GetCookieFlaky(ctx context.Context, request *proto.GetCookieFlakyRequest) (*proto.GetCookieFlakyResponse, error) {
	return api.GetCookieFlakyHandler(ctx, h.ServerContext, request)
}

func (h *Fortune) GetCookieLate(ctx context.Context, request *proto.GetCookieLateRequest) (*proto.GetCookieLateResponse, error) {
	return api.GetCookieLateHandler(ctx, h.ServerContext, request)
}
