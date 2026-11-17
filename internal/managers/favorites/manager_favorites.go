package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/helpers/response"
	commonModels "api-tests-template/internal/managers/common/models"
	"api-tests-template/internal/managers/favorites/models"
)

// maxPageSize максимальный размер страницы избранного, совпадает с лимитом избранного на пользователя
const maxPageSize = 50

// GetFavorites возвращает избранное пользователя одной страницей, ожидает 200
func GetFavorites(t *testing.T, token string) models.FavoritesList {
	resp := models.FavoritesList{}
	response.ReadJSON(t, favorites.HttpGetFavorites(t, token, maxPageSize), http.StatusOK, &resp)

	return resp
}

// GetFavoritesWithError запрашивает избранное и возвращает тело ошибки
func GetFavoritesWithError(t *testing.T, token string, expectedStatusCode int) commonModels.ErrorResponse {
	resp := commonModels.ErrorResponse{}
	response.ReadJSON(t, favorites.HttpGetFavorites(t, token, maxPageSize), expectedStatusCode, &resp)

	return resp
}

// AddFavorite добавляет объявление в избранное, ожидает 201
func AddFavorite(t *testing.T, token, advertisementId string) models.FavoriteItem {
	resp := models.FavoriteItem{}
	request := models.AddFavoriteRequest{AdvertisementId: advertisementId}
	response.ReadJSON(t, favorites.HttpPostFavorite(t, token, request), http.StatusCreated, &resp)

	return resp
}

// AddFavoriteWithError пытается добавить объявление в избранное и возвращает тело ошибки
func AddFavoriteWithError(t *testing.T, token, advertisementId string, expectedStatusCode int) commonModels.ErrorResponse {
	resp := commonModels.ErrorResponse{}
	request := models.AddFavoriteRequest{AdvertisementId: advertisementId}
	response.ReadJSON(t, favorites.HttpPostFavorite(t, token, request), expectedStatusCode, &resp)

	return resp
}

// DeleteFavorite удаляет объявление из избранного, ожидает 204
func DeleteFavorite(t *testing.T, token, advertisementId string) {
	response.ReadBody(t, favorites.HttpDeleteFavorite(t, token, advertisementId), http.StatusNoContent)
}

// DeleteFavoriteWithError пытается удалить объявление из избранного и возвращает тело ошибки
func DeleteFavoriteWithError(t *testing.T, token, advertisementId string, expectedStatusCode int) commonModels.ErrorResponse {
	resp := commonModels.ErrorResponse{}
	response.ReadJSON(t, favorites.HttpDeleteFavorite(t, token, advertisementId), expectedStatusCode, &resp)

	return resp
}
