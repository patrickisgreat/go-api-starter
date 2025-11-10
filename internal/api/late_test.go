package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/soundcloud/gokit/v2/twirp/common_session"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/patrickisgreat/go-api-starter/internal/appserver"
	"github.com/patrickisgreat/go-api-starter/internal/cookierepo"
	"github.com/patrickisgreat/go-api-starter/internal/cookieservice"
	proto "github.com/patrickisgreat/go-api-starter/internal/generated/go_api_starter/api"
)

func TestGetCookieLateHandler(t *testing.T) {
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
	request := &proto.GetCookieLateRequest{
		UserSession: &common_session.UserSession{
			UserUrn:  &wrapperspb.StringValue{Value: "myapp:users:2"},
			Features: []string{},
		},
		DelayMs: &wrapperspb.Int32Value{Value: 500},
	}

	// Capture the start time
	s := time.Now()

	// Call the handler
	response, err := GetCookieLateHandler(context.Background(), serverContext, request)

	// Assert the results
	c.Assert(err, qt.IsNil)
	c.Assert(response, qt.Not(qt.IsNil))
	c.Assert(response.FortuneCookieMessage, qt.Equals, "foo")

	f := time.Since(s)
	if f.Milliseconds() <= int64(500) {
		fmt.Printf("f.Milliseconds() = %d\n", f.Milliseconds())
	}
	c.Assert(f.Milliseconds() >= int64(500), qt.IsTrue)
}

func TestGetCookieLateHandler_InvalidSession(t *testing.T) {
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
	request := &proto.GetCookieLateRequest{
		UserSession: &common_session.UserSession{
			UserUrn:  &wrapperspb.StringValue{Value: "myapp:tracks:2"},
			Features: []string{},
		},
	}

	// Call the handler
	_, err := GetCookieLateHandler(context.Background(), serverContext, request)

	// Assert the results
	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Equals, "twirp error invalid_argument: userSession.userUrn invalid user urn")
}
