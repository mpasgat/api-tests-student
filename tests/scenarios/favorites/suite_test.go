package favorites

import (
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/suite"

	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

// createdAdvertisement объявление, созданное в тесте, и токен его владельца для удаления
type createdAdvertisement struct {
	id         string
	ownerToken string
}

type TestSuite struct {
	suite.Suite
	// loginData пользователь из .env.override, работает со своим избранным
	loginData models.LoginOkResponse
	// sellerData второй пользователь, владелец объявлений. Своё объявление добавить в избранное нельзя
	sellerData models.LoginOkResponse

	createdAdvertisements []createdAdvertisement
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Регистрация второго пользователя, который будет продавцом объявлений")
	s.sellerData = auth.RegisterRandomUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// TearDownTest удаляет объявления, созданные в тесте
func (s *TestSuite) TearDownTest() {
	advertisements := s.createdAdvertisements
	s.createdAdvertisements = nil

	for _, ad := range advertisements {
		base.Postcondition("Удаляем объявление " + ad.id)
		advManager.DeleteAdvertisement(s.T(), ad.ownerToken, ad.id)
	}
}

// createSellerAdvertisement создаёт объявление от имени продавца, после теста оно удаляется
func (s *TestSuite) createSellerAdvertisement() advModels.Advertisement {
	return s.createAdvertisement(s.sellerData.Token)
}

// createOwnAdvertisement создаёт объявление от имени основного пользователя, после теста оно удаляется
func (s *TestSuite) createOwnAdvertisement() advModels.Advertisement {
	return s.createAdvertisement(s.loginData.Token)
}

// deleteSellerAdvertisement удаляет объявление продавца в ходе теста и убирает его из списка на удаление после теста
func (s *TestSuite) deleteSellerAdvertisement(id string) {
	advManager.DeleteAdvertisement(s.T(), s.sellerData.Token, id)
	s.createdAdvertisements = slices.DeleteFunc(s.createdAdvertisements, func(ad createdAdvertisement) bool {
		return ad.id == id
	})
}

// createAdvertisement создаёт объявление от имени владельца токена и запоминает его для удаления после теста
func (s *TestSuite) createAdvertisement(token string) advModels.Advertisement {
	advertisement := advManager.CreateAdvertisement(s.T(), token, advManager.NewAdvertisementRequest())
	s.createdAdvertisements = append(s.createdAdvertisements, createdAdvertisement{id: advertisement.Id, ownerToken: token})

	return advertisement
}
