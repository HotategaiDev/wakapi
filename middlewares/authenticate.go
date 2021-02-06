package middlewares

import (
	"context"
	"fmt"
	conf "github.com/muety/wakapi/config"
	"github.com/muety/wakapi/models"
	"github.com/muety/wakapi/services"
	"github.com/muety/wakapi/utils"
	"net/http"
	"strings"
)

type AuthenticateMiddleware struct {
	config           *conf.Config
	userSrvc         services.IUserService
	optionalForPaths []string
}

func NewAuthenticateMiddleware(userService services.IUserService) *AuthenticateMiddleware {
	return &AuthenticateMiddleware{
		config:           conf.Get(),
		userSrvc:         userService,
		optionalForPaths: []string{},
	}
}

func (m *AuthenticateMiddleware) WithOptionalFor(paths []string) *AuthenticateMiddleware {
	m.optionalForPaths = paths
	return m
}

func (m *AuthenticateMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.ServeHTTP(w, r, h.ServeHTTP)
	})
}

func (m *AuthenticateMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var user *models.User
	user, err := m.tryGetUserByCookie(r)

	if err != nil {
		user, err = m.tryGetUserByApiKey(r)
	}

	if err != nil || user == nil {
		if m.isOptional(r.URL.Path) {
			next(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api") {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			http.SetCookie(w, m.config.GetClearCookie(models.AuthCookieKey, "/"))
			http.Redirect(w, r, fmt.Sprintf("%s/?error=unauthorized", m.config.Server.BasePath), http.StatusFound)
		}
		return
	}

	ctx := context.WithValue(r.Context(), models.UserKey, user)
	next(w, r.WithContext(ctx))
}

func (m *AuthenticateMiddleware) isOptional(requestPath string) bool {
	for _, p := range m.optionalForPaths {
		if strings.HasPrefix(requestPath, p) || requestPath == p {
			return true
		}
	}
	return false
}

func (m *AuthenticateMiddleware) tryGetUserByApiKey(r *http.Request) (*models.User, error) {
	key, err := utils.ExtractBearerAuth(r)
	if err != nil {
		return nil, err
	}

	var user *models.User
	userKey := strings.TrimSpace(key)
	user, err = m.userSrvc.GetUserByKey(userKey)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *AuthenticateMiddleware) tryGetUserByCookie(r *http.Request) (*models.User, error) {
	username, err := utils.ExtractCookieAuth(r, m.config)
	if err != nil {
		return nil, err
	}

	user, err := m.userSrvc.GetUserById(*username)
	if err != nil {
		return nil, err
	}

	// no need to check password here, as securecookie decoding will fail anyway,
	// if cookie is not properly signed

	return user, nil
}
