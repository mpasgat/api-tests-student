package createAdvertisement

import (
	"fmt"
	"net/http"
	"strings"

	"api-tests-template/internal/assertions"
	apiErrors "api-tests-template/internal/constants/errors"
	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"

	base "api-tests-template/tests"
)

// Ограничения сервера. В swagger указаны другие длины (255 и 1000), но фактически сервер проверяет эти
const (
	maxTitleLength       = 50
	maxDescriptionLength = 500
)

// Negative: Пытаемся создать объявление без обязательных параметров или с некорректными значениями
func (s *TestSuite) TestCreateAdvertisementInvalidParams() {
	cases := []struct {
		name            string
		modify          func(request *advModels.CreateAdvertisementRequest)
		expectedMessage string
	}{
		{
			name:            "Пустое название",
			modify:          func(r *advModels.CreateAdvertisementRequest) { r.Title = "" },
			expectedMessage: "title is required",
		},
		{
			name:   "Слишком длинное название",
			modify: func(r *advModels.CreateAdvertisementRequest) { r.Title = strings.Repeat("a", maxTitleLength+1) },
			expectedMessage: fmt.Sprintf("title must not exceed %d characters provided %d",
				maxTitleLength, maxTitleLength+1),
		},
		{
			name:            "Пустое описание",
			modify:          func(r *advModels.CreateAdvertisementRequest) { r.Description = "" },
			expectedMessage: "description is required",
		},
		{
			name: "Слишком длинное описание",
			modify: func(r *advModels.CreateAdvertisementRequest) {
				r.Description = strings.Repeat("a", maxDescriptionLength+1)
			},
			expectedMessage: fmt.Sprintf("description must be between 1 and %d characters", maxDescriptionLength),
		},
		{
			name:            "Отрицательная цена",
			modify:          func(r *advModels.CreateAdvertisementRequest) { r.Price = -1 },
			expectedMessage: "price must be between 0 and 1000000 inclusive",
		},
		{
			name:            "Нулевое количество",
			modify:          func(r *advModels.CreateAdvertisementRequest) { r.Quantity = 0 },
			expectedMessage: "quantity must be between 1 and 100 inclusive",
		},
		{
			name:            "Без фотографий",
			modify:          func(r *advModels.CreateAdvertisementRequest) { r.Photos = nil },
			expectedMessage: "at least 1 photo is required",
		},
		{
			name:            "Фотографий больше максимума",
			modify:          func(r *advModels.CreateAdvertisementRequest) { r.Photos = advManager.PhotoFilePaths() },
			expectedMessage: fmt.Sprintf("maximum %d photos allowed", advManager.MaxPhotosCount),
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			request := advManager.NewAdvertisementRequest()
			tc.modify(&request)

			base.Step("Создаём объявление: " + tc.name)
			errResp := advManager.CreateAdvertisementWithError(s.T(), s.loginData.Token, request, http.StatusBadRequest)

			assertions.ErrorResponse(s.T(), errResp, apiErrors.InvalidRequest, tc.expectedMessage)
			base.Check("Сервер вернул 400 и ошибку: " + errResp.Message)
		})
	}
}

// Negative: Пытаемся создать объявление без токена и с неверным токеном
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	cases := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{"Без токена", "", apiErrors.AuthHeaderRequiredMessage},
		{"С неверным токеном", "incorrect_token", apiErrors.InvalidAuthorizationMessage},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			base.Step("Создаём объявление: " + tc.name)
			errResp := advManager.CreateAdvertisementWithError(s.T(), tc.token, advManager.NewAdvertisementRequest(),
				http.StatusUnauthorized)

			assertions.ErrorResponse(s.T(), errResp, apiErrors.Unauthorized, tc.expectedMessage)
			base.Check("Сервер вернул 401 и ошибку: " + errResp.Message)
		})
	}
}
