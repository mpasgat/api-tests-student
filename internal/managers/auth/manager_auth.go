package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

// passwordRequiredChars добавляется к случайному паролю: сервер требует цифру, а в случайной строке её может не быть
const passwordRequiredChars = "Aa1"

// RegisterRandomUser регистрирует нового пользователя со случайным email и возвращает его токен и данные.
// Нужен, когда в тесте участвует второй пользователь, например владелец объявления для избранного
func RegisterRandomUser(t *testing.T) models.LoginOkResponse {
	request := models.RegisterRequest{
		Email:     "autotest_" + strings.ToLower(utils.RandomString(10)) + "@example.com",
		Password:  utils.RandomString(12) + passwordRequiredChars,
		FirstName: "Autotest",
	}

	registerResponse := models.LoginOkResponse{}
	response.ReadJSON(t, auth.HttpPostAuthRegister(t, request), http.StatusCreated, &registerResponse)
	require.NotEmpty(t, registerResponse.Token, "После регистрации не пришёл токен")

	return registerResponse
}
