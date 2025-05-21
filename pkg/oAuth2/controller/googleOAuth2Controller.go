package controller

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/labstack/echo/v4"
	"github.com/lekchan000/isekai-shop-api/config"
	"github.com/lekchan000/isekai-shop-api/pkg/custom"
	_oauth2Exception "github.com/lekchan000/isekai-shop-api/pkg/oAuth2/exception"
	_oauth2Model "github.com/lekchan000/isekai-shop-api/pkg/oAuth2/model"
	_oauth2Service "github.com/lekchan000/isekai-shop-api/pkg/oAuth2/service"
	"golang.org/x/oauth2"
)

type GoogleOAuth2Controller struct {
	oauth2Service _oauth2Service.OAuth2Service
	oauth2Conf    *config.OAuth2
	logger        echo.Logger
}

var (
	playerGoogleOAuth2 *oauth2.Config
	adminGoogleOAuth2  *oauth2.Config
	once               sync.Once

	accessTokenCookieName  = "act"
	refreshTokenCookieName = "rft"
	stateCookieName        = "state"

	letters = []byte("abcdefjtijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func NewGoogleOAuth2Controller(
	_oauth2Service _oauth2Service.OAuth2Service,
	oauth2Conf *config.OAuth2,
	logger echo.Logger,
) OAuth2Controller {
	once.Do(func() {
		setGoogleOAuth2Config(oauth2Conf)
	})

	return &GoogleOAuth2Controller{
		_oauth2Service,
		oauth2Conf,
		logger,
	}
}

func setGoogleOAuth2Config(oauth2Conf *config.OAuth2) {
	playerGoogleOAuth2 = &oauth2.Config{
		ClientID:     oauth2Conf.ClientID,
		ClientSecret: oauth2Conf.ClientSecret,
		RedirectURL:  oauth2Conf.PlayerRedirectUrl,
		Scopes:       oauth2Conf.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:       oauth2Conf.EndPoints.AuthUrl,
			TokenURL:      oauth2Conf.EndPoints.TokenUrl,
			DeviceAuthURL: oauth2Conf.EndPoints.DeviceAuthUrl,
			AuthStyle:     oauth2.AuthStyleInParams,
		},
	}

	adminGoogleOAuth2 = &oauth2.Config{
		ClientID:     oauth2Conf.ClientID,
		ClientSecret: oauth2Conf.ClientSecret,
		RedirectURL:  oauth2Conf.AdminRedirectUrl,
		Scopes:       oauth2Conf.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:       oauth2Conf.EndPoints.AuthUrl,
			TokenURL:      oauth2Conf.EndPoints.TokenUrl,
			DeviceAuthURL: oauth2Conf.EndPoints.DeviceAuthUrl,
			AuthStyle:     oauth2.AuthStyleInParams,
		},
	}
}

func (c *GoogleOAuth2Controller) Playerlogin(pctx echo.Context) error {
	state := c.randomState()
	c.setCookie(pctx, stateCookieName, state)

	return pctx.Redirect(http.StatusFound, playerGoogleOAuth2.AuthCodeURL(state))
}

func (c *GoogleOAuth2Controller) Adminlogin(pctx echo.Context) error {
	state := c.randomState()
	c.setCookie(pctx, stateCookieName, state)

	return pctx.Redirect(http.StatusFound, adminGoogleOAuth2.AuthCodeURL(state))
}

func (c *GoogleOAuth2Controller) PlayerloginCallback(pctx echo.Context) error {
	//ctx := context.Background()
	if err := retry.Do(func() error {
		return c.callbackValidating(pctx)
	}, retry.Attempts(3), retry.Delay(3*time.Second)); err != nil {
		c.logger.Errorf("Failed to validate callback: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, err.Error())
	}
	return pctx.JSON(http.StatusOK, _oauth2Model.LoginResponse{Message: "Login Success"})
}

func (c *GoogleOAuth2Controller) AdminloginCallback(pctx echo.Context) error {
	//ctx := context.Background()
	if err := retry.Do(func() error {
		return c.callbackValidating(pctx)
	}, retry.Attempts(3), retry.Delay(3*time.Second)); err != nil {
		c.logger.Errorf("Failed to validate callback: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, err.Error())
	}
	return pctx.JSON(http.StatusOK, _oauth2Model.LoginResponse{Message: "Login Success"})
}

func (c *GoogleOAuth2Controller) Logout(pctx echo.Context) error {
	c.removeCookie(pctx, accessTokenCookieName)
	c.removeCookie(pctx, refreshTokenCookieName)
	c.removeCookie(pctx, stateCookieName)

	return pctx.NoContent(http.StatusNoContent)
}

func (c *GoogleOAuth2Controller) setCookie(pctx echo.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	pctx.SetCookie(cookie)
}

func (c *GoogleOAuth2Controller) removeCookie(pctx echo.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	pctx.SetCookie(cookie)
}

func (c *GoogleOAuth2Controller) callbackValidating(pctx echo.Context) error {
	state := pctx.QueryParam("state")
	stateFromCookie, err := pctx.Cookie(stateCookieName)
	if err != nil {
		c.logger.Errorf("Failed to get state from cookie: %v", err.Error())
		return &_oauth2Exception.Unauthorized{}
	}

	if state != stateFromCookie.Value {
		c.logger.Errorf("Invalid state: %s", state)
		return &_oauth2Exception.Unauthorized{}
	}
	return nil
}

func (c *GoogleOAuth2Controller) randomState() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
