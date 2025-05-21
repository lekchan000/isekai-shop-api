package server

import (
	_adminRepository "github.com/lekchan000/isekai-shop-api/pkg/admin/repository"
	_oauth2Controller "github.com/lekchan000/isekai-shop-api/pkg/oAuth2/controller"
	_oauth2Service "github.com/lekchan000/isekai-shop-api/pkg/oAuth2/service"
	_playerRepository "github.com/lekchan000/isekai-shop-api/pkg/player/repository"
)

func (s *echoServer) initOAuth2Router() {
	router := s.app.Group("v1/oauth2/google")

	playerRepository := _playerRepository.NewPlayerRepositoryImpl(s.db, s.app.Logger)
	adminRepository := _adminRepository.NewAdminRepositoryImpl(s.db, s.app.Logger)

	oauth2Service := _oauth2Service.NewGoogleOAuth2Service(playerRepository, adminRepository)
	oauth2Controller := _oauth2Controller.NewGoogleOAuth2Controller(oauth2Service, s.conf.OAuth2, s.app.Logger)

	router.GET("/player/login", oauth2Controller.Playerlogin)
	router.GET("/admin/login", oauth2Controller.Adminlogin)
	router.GET("/player/login/callback", oauth2Controller.PlayerloginCallback)
	router.GET("/admin/login/callback", oauth2Controller.AdminloginCallback)
	router.DELETE("/logout", oauth2Controller.Logout)
}
