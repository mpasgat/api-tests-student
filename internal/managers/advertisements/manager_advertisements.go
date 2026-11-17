package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements/models"
	commonModels "api-tests-template/internal/managers/common/models"
	"api-tests-template/internal/utils"
)

const (
	// MaxPhotosCount максимальное количество фотографий в объявлении
	MaxPhotosCount = 3

	defaultDescription = "Объявление создано автотестом"
	defaultPrice       = 1499.99
	defaultQuantity    = 3
)

// photoFilePaths тестовые фотографии относительно корня проекта.
// Их на одну больше, чем MaxPhotosCount, чтобы проверить превышение лимита
var photoFilePaths = []string{
	"testdata/photo.jpg",
	"testdata/photo2.png",
	"testdata/photo3.png",
	"testdata/photo4.png",
}

// PhotoFilePaths возвращает абсолютные пути ко всем тестовым фотографиям
func PhotoFilePaths() []string {
	paths := make([]string, 0, len(photoFilePaths))
	for _, photo := range photoFilePaths {
		paths = append(paths, utils.ProjectFilePath(photo))
	}

	return paths
}

// NewAdvertisementRequest возвращает корректный запрос на создание объявления со всеми параметрами
// и максимальным количеством фотографий. Название случайное, чтобы объявления из разных прогонов не пересекались
func NewAdvertisementRequest() models.CreateAdvertisementRequest {
	return models.CreateAdvertisementRequest{
		Title:       "Autotest " + utils.RandomString(10),
		Description: defaultDescription,
		Price:       defaultPrice,
		Quantity:    defaultQuantity,
		Photos:      PhotoFilePaths()[:MaxPhotosCount],
	}
}

// GetAdvertisementsList возвращает список объявлений, найденных по подстроке в названии
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := advertisements.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// CreateAdvertisement создаёт объявление и возвращает его, ожидает 201
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	resp := models.Advertisement{}
	response.ReadJSON(t, advertisements.HttpPostAdvertisement(t, token, request), http.StatusCreated, &resp)
	require.NotEmpty(t, resp.Id, "У созданного объявления пустой id")

	return resp
}

// CreateAdvertisementWithError пытается создать объявление с некорректными данными и возвращает тело ошибки
func CreateAdvertisementWithError(t *testing.T, token string, request models.CreateAdvertisementRequest,
	expectedStatusCode int) commonModels.ErrorResponse {
	resp := commonModels.ErrorResponse{}
	response.ReadJSON(t, advertisements.HttpPostAdvertisement(t, token, request), expectedStatusCode, &resp)

	return resp
}

// GetAdvertisement возвращает объявление по id, ожидает 200
func GetAdvertisement(t *testing.T, id string) models.Advertisement {
	resp := models.Advertisement{}
	response.ReadJSON(t, advertisements.HttpGetAdvertisement(t, id), http.StatusOK, &resp)

	return resp
}

// GetAdvertisementPhotos возвращает список фотографий объявления, ожидает 200
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	var resp []models.Photo
	response.ReadJSON(t, advertisements.HttpGetAdvertisementPhotos(t, token, id), http.StatusOK, &resp)

	return resp
}

// DeleteAdvertisement удаляет объявление, ожидает 204
func DeleteAdvertisement(t *testing.T, token, id string) {
	response.ReadBody(t, advertisements.HttpDeleteAdvertisement(t, token, id), http.StatusNoContent)
}
