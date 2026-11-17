package favorites

import (
	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	favManager "api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// Positive: Добавляем объявление в избранное, проверяем, что оно в списке, удаляем и проверяем, что пропало
func (s *TestSuite) TestAddAndRemoveFavorite() {
	base.Precondition("Продавец создаёт объявление")
	advertisement := s.createSellerAdvertisement()

	s.Run("Добавляем объявление в избранное", func() {
		base.Step("Добавляем в избранное объявление " + advertisement.Id)
		favorite := favManager.AddFavorite(s.T(), s.loginData.Token, advertisement.Id)

		require.Equal(s.T(), advertisement.Id, favorite.Advertisement.Id, "В ответе другое объявление")
		require.NotEmpty(s.T(), favorite.FavoriteId, "Пустой favorite_id")
		base.Check("Объявление добавлено в избранное")
	})

	s.Run("Проверяем, что объявление появилось в избранном", func() {
		base.Step("Получаем список избранного")
		assertions.ContainsFavorite(s.T(), favManager.GetFavorites(s.T(), s.loginData.Token), advertisement.Id)
		base.Check("Объявление есть в избранном")
	})

	s.Run("Удаляем объявление из избранного", func() {
		base.Step("Удаляем из избранного объявление " + advertisement.Id)
		favManager.DeleteFavorite(s.T(), s.loginData.Token, advertisement.Id)
		base.Check("Объявление удалено из избранного")
	})

	s.Run("Проверяем, что объявление пропало из избранного", func() {
		base.Step("Получаем список избранного")
		assertions.NotContainsFavorite(s.T(), favManager.GetFavorites(s.T(), s.loginData.Token), advertisement.Id)
		base.Check("Объявления больше нет в избранном")
	})
}

// Positive: Удаляем объявление, которое лежит в избранном, и проверяем, что оно пропало из избранного
func (s *TestSuite) TestDeletedAdvertisementDisappearsFromFavorites() {
	base.Precondition("Продавец создаёт объявление, пользователь добавляет его в избранное")
	advertisement := s.createSellerAdvertisement()
	favManager.AddFavorite(s.T(), s.loginData.Token, advertisement.Id)

	s.Run("Проверяем, что объявление есть в избранном", func() {
		base.Step("Получаем список избранного")
		assertions.ContainsFavorite(s.T(), favManager.GetFavorites(s.T(), s.loginData.Token), advertisement.Id)
		base.Check("Объявление есть в избранном")
	})

	s.Run("Продавец удаляет объявление", func() {
		base.Step("Удаляем объявление " + advertisement.Id)
		s.deleteSellerAdvertisement(advertisement.Id)
		base.Check("Объявление удалено")
	})

	s.Run("Проверяем, что объявление пропало из избранного", func() {
		base.Step("Получаем список избранного")
		assertions.NotContainsFavorite(s.T(), favManager.GetFavorites(s.T(), s.loginData.Token), advertisement.Id)
		base.Check("Удалённого объявления нет в избранном")
	})
}
