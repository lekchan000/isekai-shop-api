package controller

import "github.com/labstack/echo/v4"

type ItemManagingController interface {
	Creating(ptcx echo.Context) error
	Editing(pctx echo.Context) error
	Achiving(pctx echo.Context) error
}
