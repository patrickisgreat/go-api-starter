package api

import (
	"slices"

	common "github.com/soundcloud/gokit/v2/twirp/common_session"
	"github.com/soundcloud/gokit/v2/urn"
	"github.com/twitchtv/twirp"
)

func isValidUrn(urnValue string, urnTypes []string) bool {
	if urnValue == "" {
		return false
	}

	u, err := urn.Parse(urnValue)
	if err != nil || (len(urnTypes) > 0 && !slices.Contains(urnTypes, u.Collection)) {
		return false
	}

	return true
}

func validateUserSession(session *common.UserSession) twirp.Error {
	if session.GetUserUrn().GetValue() != "" && !isValidUrn(session.GetUserUrn().GetValue(), []string{"users"}) {
		return twirp.InvalidArgumentError("userSession.userUrn", "invalid user urn")
	}

	return nil
}
