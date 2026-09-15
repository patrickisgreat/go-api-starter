package end_to_end

import (
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestGetQuote(t *testing.T) {
	c := qt.New(t)

	url := apiHost + "/twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/GetQuote"
	payload := map[string]interface{}{
		"user_session": map[string]interface{}{
			"user_urn": "myapp:users:2",
			"features": []string{},
		},
	}
	response, status, err := postRequest(url, payload)
	c.Assert(err, qt.IsNil)
	c.Assert(status, qt.Equals, 200)
	c.Assert(response["quote"], qt.IsNotNil)
}

func TestGetQuoteFlaky(t *testing.T) {
	c := qt.New(t)

	url := apiHost + "/twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/GetQuoteFlaky"
	payload := map[string]interface{}{
		"user_session": map[string]interface{}{
			"user_urn": "myapp:users:2",
			"features": []string{},
		},
		"fail": true,
	}
	response, status, err := postRequest(url, payload)
	c.Assert(err, qt.IsNil)
	c.Assert(status, qt.Equals, 500)
	c.Assert(response["msg"], qt.Equals, "flaky error")
}

func TestGetQuoteLate(t *testing.T) {
	c := qt.New(t)

	url := apiHost + "/twirp/proto.patrickisgreat.pb_go_api_starter.api.Quotes/GetQuoteLate"
	payload := map[string]interface{}{
		"user_session": map[string]interface{}{
			"user_urn": "myapp:users:2",
			"features": []string{},
		},
		"delay_ms": 500,
	}
	s := time.Now()
	response, status, err := postRequest(url, payload)
	c.Assert(err, qt.IsNil)
	c.Assert(status, qt.Equals, 200)
	c.Assert(response["quote"], qt.IsNotNil)
	f := time.Since(s)
	c.Assert(f.Milliseconds() > int64(500), qt.IsTrue)
}
