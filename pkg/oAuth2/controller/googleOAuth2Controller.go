package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/labstack/echo/v4"
	"github.com/lekchan000/isekai-shop-api/config"
	_adminModel "github.com/lekchan000/isekai-shop-api/pkg/admin/model"
	"github.com/lekchan000/isekai-shop-api/pkg/custom"
	_oauth2Exception "github.com/lekchan000/isekai-shop-api/pkg/oAuth2/exception"
	_oauth2Model "github.com/lekchan000/isekai-shop-api/pkg/oAuth2/model"
	_oauth2Service "github.com/lekchan000/isekai-shop-api/pkg/oAuth2/service"
	_playerModel "github.com/lekchan000/isekai-shop-api/pkg/player/model"
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
	ctx := context.Background()
	if err := retry.Do(func() error {
		return c.callbackValidating(pctx)
	}, retry.Attempts(3), retry.Delay(3*time.Second)); err != nil {
		c.logger.Errorf("Failed to validate callback: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, err)
	}

	token, err := playerGoogleOAuth2.Exchange(ctx, pctx.QueryParam("code"))
	if err != nil {
		c.logger.Errorf("Failed to exchange token: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	client := playerGoogleOAuth2.Client(ctx, token)

	userInfo, err := c.getUserInfo(client)
	if err != nil {
		c.logger.Errorf("Failed to get user info: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	playerCreatingReq := &_playerModel.PlayerCreatingReq{
		ID:     userInfo.ID,
		Email:  userInfo.Email,
		Name:   userInfo.Name,
		Avatar: userInfo.Picture,
	}

	if err := c.oauth2Service.PlayerAccountCreating(playerCreatingReq); err != nil {
		c.logger.Errorf("Failed to create account: %s", err.Error())
		return custom.Error(pctx, http.StatusInternalServerError, &_oauth2Exception.OAuth2Processing{})
	}

	c.setSameSiteCookie(pctx, accessTokenCookieName, token.AccessToken)
	c.setSameSiteCookie(pctx, refreshTokenCookieName, token.RefreshToken)

	return pctx.JSON(http.StatusOK, _oauth2Model.LoginResponse{Message: "Login Success"})
}

func (c *GoogleOAuth2Controller) AdminloginCallback(pctx echo.Context) error {
	ctx := context.Background()
	if err := retry.Do(func() error {
		return c.callbackValidating(pctx)
	}, retry.Attempts(3), retry.Delay(3*time.Second)); err != nil {
		c.logger.Errorf("Failed to validate callback: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, err)
	}

	token, err := adminGoogleOAuth2.Exchange(ctx, pctx.QueryParam("code"))
	if err != nil {
		c.logger.Errorf("Failed to exchange token: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	client := adminGoogleOAuth2.Client(ctx, token)

	userInfo, err := c.getUserInfo(client)
	if err != nil {
		c.logger.Errorf("Failed to get user info: %s", err.Error())
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	adminCreatingReq := &_adminModel.AdminCreatingReq{
		ID:     userInfo.ID,
		Email:  userInfo.Email,
		Name:   userInfo.Name,
		Avatar: userInfo.Picture,
	}

	if err := c.oauth2Service.AdminAccountCreating(adminCreatingReq); err != nil {
		c.logger.Errorf("Failed to create account: %s", err.Error())
		return custom.Error(pctx, http.StatusInternalServerError, &_oauth2Exception.OAuth2Processing{})
	}

	c.setSameSiteCookie(pctx, accessTokenCookieName, token.AccessToken)
	c.setSameSiteCookie(pctx, refreshTokenCookieName, token.RefreshToken)

	return pctx.JSON(http.StatusOK, _oauth2Model.LoginResponse{Message: "Login Success"})
}

func (c *GoogleOAuth2Controller) Logout(pctx echo.Context) error {
	accessToken, err := pctx.Cookie(accessTokenCookieName)
	if err != nil {
		c.logger.Errorf("error reading access token: %s", err.Error())
		return custom.Error(pctx, http.StatusInternalServerError, &_oauth2Exception.Logout{})
	}

	if err := c.revokeToken(accessToken.Value); err != nil {
		c.logger.Errorf("error revoking token: %s", err.Error())
		return custom.Error(pctx, http.StatusInternalServerError, &_oauth2Exception.Logout{})
	}

	c.removeSameSiteCookie(pctx, accessTokenCookieName)
	c.removeSameSiteCookie(pctx, refreshTokenCookieName)

	return pctx.JSON(http.StatusOK, &_oauth2Model.LoginResponse{Message: "Logout Success"})
}

func (c *GoogleOAuth2Controller) revokeToken(accessToken string) error {
	revokeURL := fmt.Sprintf("%s?token=%s", c.oauth2Conf.RevokeUrl, accessToken)

	resp, err := http.Post(revokeURL, "application/x-www-from-urlencoded", nil)
	if err != nil {
		fmt.Println("Error revoking token: ", err)
		return err
	}

	defer resp.Body.Close()

	return nil
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

func (c *GoogleOAuth2Controller) setSameSiteCookie(pctx echo.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	pctx.SetCookie(cookie)
}

func (c *GoogleOAuth2Controller) removeSameSiteCookie(pctx echo.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	}
	pctx.SetCookie(cookie)
}

func (c *GoogleOAuth2Controller) getUserInfo(client *http.Client) (*_oauth2Model.UserInfo, error) {
	resp, err := client.Get(c.oauth2Conf.UserInfoUrl)
	if err != nil {
		c.logger.Errorf("error getting user info: %s", err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	userInfoInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Errorf("error reading user info: %s", err.Error())
		return nil, err
	}

	userinfo := new(_oauth2Model.UserInfo)
	if err := json.Unmarshal(userInfoInBytes, &userinfo); err != nil {
		c.logger.Errorf("error unmarshalling user info: %s", err.Error())
		return nil, err
	}
	return userinfo, nil
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

	c.removeCookie(pctx, stateCookieName)

	return nil
}

func (c *GoogleOAuth2Controller) randomState() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
