package createAdvertisement

import (
	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"

	base "api-tests-template/tests"
)

// Positive: Создаём объявление со всеми параметрами и проверяем, что оно доступно по id, в поиске и с фотографиями
func (s *TestSuite) TestCreateAdvertisementPositive() {
	request := advManager.NewAdvertisementRequest()

	var created advModels.Advertisement
	s.Run("Создаём объявление со всеми параметрами", func() {
		base.Step("Создаём объявление с названием " + request.Title)
		created = s.createAdvertisement(request)
		base.Check("Объявление создано, id " + created.Id)
	})

	s.Run("Проверяем, что объявление по id возвращает переданные параметры", func() {
		base.Step("Получаем объявление по id")
		advertisement := advManager.GetAdvertisement(s.T(), created.Id)

		require.Equal(s.T(), created.Id, advertisement.Id, "id не совпадает")
		require.Equal(s.T(), s.loginData.User.Id, advertisement.UserId, "user_id не совпадает")
		require.Equal(s.T(), request.Title, advertisement.Title, "Название не совпадает")
		require.Equal(s.T(), request.Description, advertisement.Description, "Описание не совпадает")
		require.Equal(s.T(), request.Price, advertisement.Price, "Цена не совпадает")
		require.Equal(s.T(), request.Quantity, advertisement.Quantity, "Количество не совпадает")
		base.Check("Параметры объявления совпадают с переданными при создании")
	})

	s.Run("Проверяем, что объявление находится в поиске по полному названию", func() {
		base.Step("Ищем объявление по названию " + request.Title)
		list := advManager.GetAdvertisementsList(s.T(), s.loginData.Token, request.Title)

		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, request.Title)
		base.Check("Объявление найдено в поиске")
	})

	s.Run("Проверяем, что фотографии загружены и доступны", func() {
		base.Step("Получаем фотографии объявления")
		photos := advManager.GetAdvertisementPhotos(s.T(), s.loginData.Token, created.Id)

		require.Len(s.T(), photos, len(request.Photos), "Количество фотографий не совпадает с загруженным")
		require.ElementsMatch(s.T(), photoIds(created.Photos), photoIds(photos),
			"Фотографии объявления не совпадают с загруженными при создании")
		for _, photo := range photos {
			require.NotEmptyf(s.T(), photo.Url, "У фотографии %s пустой url", photo.Id)
			base.Check("Фотография загружена, url " + photo.Url)
		}
	})
}

// photoIds собирает id фотографий, чтобы сравнить наборы без учёта порядка
func photoIds(photos []advModels.Photo) []string {
	ids := make([]string, 0, len(photos))
	for _, photo := range photos {
		ids = append(ids, photo.Id)
	}

	return ids
}
