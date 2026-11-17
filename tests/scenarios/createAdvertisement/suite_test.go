package createAdvertisement

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse

	// createdAdvertisementIds объявления, созданные в текущем тесте, удаляются в TearDownTest
	createdAdvertisementIds []string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// TearDownTest удаляет объявления, созданные в тесте
func (s *TestSuite) TearDownTest() {
	ids := s.createdAdvertisementIds
	s.createdAdvertisementIds = nil

	for _, id := range ids {
		base.Postcondition("Удаляем объявление " + id)
		advManager.DeleteAdvertisement(s.T(), s.loginData.Token, id)
	}
}

// createAdvertisement создаёт объявление и запоминает его id, чтобы удалить после теста
func (s *TestSuite) createAdvertisement(request advModels.CreateAdvertisementRequest) advModels.Advertisement {
	advertisement := advManager.CreateAdvertisement(s.T(), s.loginData.Token, request)
	s.createdAdvertisementIds = append(s.createdAdvertisementIds, advertisement.Id)

	return advertisement
}
