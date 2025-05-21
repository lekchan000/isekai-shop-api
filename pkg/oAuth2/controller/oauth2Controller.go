package controller

import "github.com/labstack/echo/v4"

type OAuth2Controller interface {
	Playerlogin(pctx echo.Context) error
	Adminlogin(pctx echo.Context) error
	PlayerloginCallback(pctx echo.Context) error
	AdminloginCallback(pctx echo.Context) error
	Logout(pctx echo.Context) error
}
