package auth

import (
	"net/http"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/db"
	"github.com/concourse/skymarshal/provider"
	"github.com/concourse/skymarshal/routes"
	"github.com/dgrijalva/jwt-go"
	"github.com/tedsuo/rata"
)

var SigningMethod = jwt.SigningMethodRS256

//go:generate counterfeiter . ProviderFactory

type ProviderFactory interface {
	GetProvider(db.Team, string) (provider.Provider, bool, error)
}

func NewOAuthHandler(
	logger lager.Logger,
	providerFactory ProviderFactory,
	teamFactory db.TeamFactory,
	csrfTokenGenerator CSRFTokenGenerator,
	authTokenGenerator AuthTokenGenerator,
	expire time.Duration,
	isTLSEnabled bool,
) (http.Handler, error) {
	return rata.NewRouter(
		routes.OAuthRoutes,
		map[string]http.Handler{
			routes.OAuthBegin: NewOAuthBeginHandler(
				logger.Session("oauth-begin"),
				providerFactory,
				teamFactory,
				expire,
				isTLSEnabled,
			),
			routes.OAuthCallback: NewOAuthCallbackHandler(
				logger.Session("oauth-callback"),
				providerFactory,
				teamFactory,
				csrfTokenGenerator,
				authTokenGenerator,
				expire,
				isTLSEnabled,
				oauthV2StateValidator{},
			),
			routes.LogOut: NewLogOutHandler(
				logger.Session("logout"),
			),
		},
	)
}

func NewOAuthV1Handler(
	logger lager.Logger,
	providerFactory ProviderFactory,
	teamFactory db.TeamFactory,
	csrfTokenGenerator CSRFTokenGenerator,
	authTokenGenerator AuthTokenGenerator,
	expire time.Duration,
	isTLSEnabled bool,
) (http.Handler, error) {
	return rata.NewRouter(
		routes.OAuthV1Routes,
		map[string]http.Handler{
			routes.OAuthV1Begin: NewOAuthBeginHandler(
				logger.Session("oauth-v1-begin"),
				providerFactory,
				teamFactory,
				expire,
				isTLSEnabled,
			),
			routes.OAuthV1Callback: NewOAuthCallbackHandler(
				logger.Session("oauth-v1-callback"),
				providerFactory,
				teamFactory,
				csrfTokenGenerator,
				authTokenGenerator,
				expire,
				isTLSEnabled,
				oauthV1StateValidator{},
			),
		},
	)
}
