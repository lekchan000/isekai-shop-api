package controller

import "github.com/labstack/echo/v4"

type ItemManagingController interface {
	Creating(ptcx echo.Context) error
}
