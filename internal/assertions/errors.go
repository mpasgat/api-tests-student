package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	commonModels "api-tests-template/internal/managers/common/models"
)

// ErrorResponse проверяет код и текст ошибки в ответе сервера
func ErrorResponse(t *testing.T, resp commonModels.ErrorResponse, expectedError, expectedMessage string) {
	require.Equal(t, expectedError, resp.Error, "Некорректное значение error параметра")
	require.Equal(t, expectedMessage, resp.Message, "Некорректное значение message параметра")
}
