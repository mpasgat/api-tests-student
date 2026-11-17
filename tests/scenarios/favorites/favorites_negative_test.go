package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	apiErrors "api-tests-template/internal/constants/errors"
	favManager "api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// nonExistentAdvertisementId корректный uuid, которого нет среди объявлений.
// Нулевой uuid не подходит: сервер считает его пустым значением и отвечает 400
const nonExistentAdvertisementId = "11111111-2222-3333-4444-555555555555"

// Negative: Добавляем в избранное и удаляем из избранного несуществующее объявление
func (s *TestSuite) TestFavoritesNonExistentAdvertisement() {
	s.Run("Добавляем несуществующее объявление в избранное", func() {
		base.Step("Добавляем в избранное объявление " + nonExistentAdvertisementId)
		errResp := favManager.AddFavoriteWithError(s.T(), s.loginData.Token, nonExistentAdvertisementId,
			http.StatusNotFound)

		assertions.ErrorResponse(s.T(), errResp, apiErrors.NotFound, "Advertisement not found")
		base.Check("Сервер вернул 404 и ошибку: " + errResp.Message)
	})

	s.Run("Удаляем из избранного несуществующее объявление", func() {
		base.Step("Удаляем из избранного объявление " + nonExistentAdvertisementId)
		errResp := favManager.DeleteFavoriteWithError(s.T(), s.loginData.Token, nonExistentAdvertisementId,
			http.StatusNotFound)

		assertions.ErrorResponse(s.T(), errResp, apiErrors.NotFound, "Favorite entry not found")
		base.Check("Сервер вернул 404 и ошибку: " + errResp.Message)
	})
}

// Negative: Вызываем ручки избранного без токена и с неверным токеном
func (s *TestSuite) TestFavoritesUnauthorized() {
	tokens := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{"без токена", "", apiErrors.AuthHeaderRequiredMessage},
		{"с неверным токеном", "incorrect_token", apiErrors.InvalidAuthorizationMessage},
	}

	for _, tc := range tokens {
		s.Run("Получаем избранное "+tc.name, func() {
			base.Step("Получаем избранное " + tc.name)
			errResp := favManager.GetFavoritesWithError(s.T(), tc.token, http.StatusUnauthorized)
			assertions.ErrorResponse(s.T(), errResp, apiErrors.Unauthorized, tc.expectedMessage)
			base.Check("GET /favorites вернул 401 и ошибку: " + errResp.Message)
		})

		s.Run("Добавляем в избранное "+tc.name, func() {
			base.Step("Добавляем в избранное " + tc.name)
			errResp := favManager.AddFavoriteWithError(s.T(), tc.token, nonExistentAdvertisementId,
				http.StatusUnauthorized)
			assertions.ErrorResponse(s.T(), errResp, apiErrors.Unauthorized, tc.expectedMessage)
			base.Check("POST /favorites вернул 401 и ошибку: " + errResp.Message)
		})

		s.Run("Удаляем из избранного "+tc.name, func() {
			base.Step("Удаляем из избранного " + tc.name)
			errResp := favManager.DeleteFavoriteWithError(s.T(), tc.token, nonExistentAdvertisementId,
				http.StatusUnauthorized)
			assertions.ErrorResponse(s.T(), errResp, apiErrors.Unauthorized, tc.expectedMessage)
			base.Check("DELETE /favorites/{ad_id} вернул 401 и ошибку: " + errResp.Message)
		})
	}
}

// Negative: Повторно добавляем в избранное то же объявление
func (s *TestSuite) TestAddFavoriteTwice() {
	base.Precondition("Продавец создаёт объявление, пользователь добавляет его в избранное")
	advertisement := s.createSellerAdvertisement()
	favManager.AddFavorite(s.T(), s.loginData.Token, advertisement.Id)

	base.Step("Добавляем то же объявление в избранное ещё раз")
	errResp := favManager.AddFavoriteWithError(s.T(), s.loginData.Token, advertisement.Id, http.StatusConflict)

	assertions.ErrorResponse(s.T(), errResp, apiErrors.AlreadyExists, "Advertisement is already in favorites")
	base.Check("Сервер вернул 409 и ошибку: " + errResp.Message)
}

// Negative: Добавляем в избранное своё собственное объявление
func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	base.Precondition("Пользователь создаёт своё объявление")
	advertisement := s.createOwnAdvertisement()

	base.Step("Добавляем своё объявление в избранное")
	errResp := favManager.AddFavoriteWithError(s.T(), s.loginData.Token, advertisement.Id, http.StatusForbidden)

	assertions.ErrorResponse(s.T(), errResp, apiErrors.Forbidden, "Cannot add own advertisement to favorites")
	base.Check("Сервер вернул 403 и ошибку: " + errResp.Message)
}

// Negative: Добавляем в избранное с пустым и некорректным advertisement_id
func (s *TestSuite) TestAddFavoriteInvalidAdvertisementId() {
	cases := []struct {
		name            string
		advertisementId string
		expectedMessage string
	}{
		{"Пустой advertisement_id", "", "invalid UUID length: 0"},
		{"advertisement_id не в формате uuid", "not-a-uuid", "invalid UUID length: 10"},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			base.Step("Добавляем в избранное: " + tc.name)
			errResp := favManager.AddFavoriteWithError(s.T(), s.loginData.Token, tc.advertisementId,
				http.StatusBadRequest)

			assertions.ErrorResponse(s.T(), errResp, apiErrors.InvalidRequest, tc.expectedMessage)
			base.Check("Сервер вернул 400 и ошибку: " + errResp.Message)
		})
	}
}
