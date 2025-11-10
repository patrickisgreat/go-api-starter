package cookierepo

import (
	"context"
	"embed"
	"encoding/json"
	"log/slog"
	"math/rand"

	"github.com/soundcloud/gokit/v2/logger"
)

//go:embed fortunes.json
var f embed.FS

type CookieRepository struct {
	Fortunes []string `json:"fortunes"`
}

func NewCookieRepository() *CookieRepository {
	fortunes := &CookieRepository{}
	data, err := f.ReadFile("fortunes.json")
	if err == nil {
		_ = json.Unmarshal(data, fortunes)
	}
	return fortunes
}

func (cr *CookieRepository) Gimme(ctx context.Context) (string, error) {
	log := logger.Get(ctx)
	cookieCount := len(cr.Fortunes)
	if cookieCount == 0 {
		err := &NoCookiesError{}
		log.Error("I cant get no cookies", slog.Any("err", err))
		return "", err
	}
	ix := rand.Intn(len(cr.Fortunes))
	return cr.Fortunes[ix], nil
}
