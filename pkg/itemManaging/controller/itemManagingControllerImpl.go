package controller

import (
	_itemManagingService "github.com/lekchan000/isekai-shop-api/pkg/itemManaging/service"
)

type itemManagingControllerImpl struct {
	itemManagingService _itemManagingService.ItemManagingService
}

func NewItemManagingControllerImpl(itemManagingService _itemManagingService.ItemManagingService) ItemManagingController {
	return &itemManagingControllerImpl{itemManagingService}
}
