package api

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/soundcloud/gokit/v2/twirp/common_session"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/patrickisgreat/go-api-starter/internal/appserver"
	"github.com/patrickisgreat/go-api-starter/internal/cookierepo"
	"github.com/patrickisgreat/go-api-starter/internal/cookieservice"
	proto "github.com/patrickisgreat/go-api-starter/internal/generated/go_api_starter/api"
)

func TestGetCookieFlakyHandlerFail(t *testing.T) {
	c := qt.New(t)

	// Mock the CookieService
	mockService := cookieservice.NewCookieService(&cookierepo.CookieRepository{
		Fortunes: []string{"foo"},
	})

	// Create a mock ServerContext
	serverContext := &appserver.ServerContext{
		CookieService: mockService,
	}

	// Create a request
	request := &proto.GetCookieFlakyRequest{
		UserSession: &common_session.UserSession{
			UserUrn:  &wrapperspb.StringValue{Value: "myapp:users:2"},
			Features: []string{},
		},
		Fail: true,
	}

	// Call the handler
	_, err := GetCookieFlakyHandler(context.Background(), serverContext, request)

	// Assert the results
	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Equals, "twirp error internal: flaky error")
}

func TestGetCookieFlakyHandler_InvalidSession(t *testing.T) {
	c := qt.New(t)

	// Mock the CookieService
	mockService := cookieservice.NewCookieService(&cookierepo.CookieRepository{
		Fortunes: []string{"foo"},
	})

	// Create a mock ServerContext
	serverContext := &appserver.ServerContext{
		CookieService: mockService,
	}

	// Create a request
	request := &proto.GetCookieFlakyRequest{
		UserSession: &common_session.UserSession{
			UserUrn:  &wrapperspb.StringValue{Value: "myapp:tracks:2"},
			Features: []string{},
		},
	}

	// Call the handler
	_, err := GetCookieFlakyHandler(context.Background(), serverContext, request)

	// Assert the results
	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Equals, "twirp error invalid_argument: userSession.userUrn invalid user urn")
}
