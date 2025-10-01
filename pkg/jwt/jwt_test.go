package jwt

import (
	"testing"
	"time"

	"log"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	// Загружаем переменные из .env файла
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func TestGenerateTokenPair(t *testing.T) {
	const email = "test@example.com"
	jwtService := NewJWT("test-secret")

	// Устанавливаем короткое время жизни для тестов
	jwtService.AccessTokenTTL = time.Second * 10
	jwtService.RefreshTokenTTL = time.Second * 20

	tokenPair, err := jwtService.GenerateTokenPair(JWTData{
		Email: email,
	})

	require.NoError(t, err)

	if tokenPair.AccessToken == "" {
		t.Fatal("Access token is empty")
	}

	if tokenPair.RefreshToken == "" {
		t.Fatal("Refresh token is empty")
	}
}

func TestParseAccessToken(t *testing.T) {
	const email = "test@example.com"
	jwtService := NewJWT("test-secret")

	tokenPair, err := jwtService.GenerateTokenPair(JWTData{
		Email: email,
	})

	require.NoError(t, err)

	isValid, data := jwtService.ParseToken(tokenPair.AccessToken)

	require.True(t, isValid)
	require.Equal(t, email, data.Email)

}

func TestParseRefreshToken(t *testing.T) {
	const email = "test@example.com"
	jwtService := NewJWT("test-secret")

	tokenPair, err := jwtService.GenerateTokenPair(JWTData{
		Email: email,
	})

	require.NoError(t, err)

	isValid, data := jwtService.ParseRefreshToken(tokenPair.RefreshToken)

	require.True(t, isValid)
	require.Equal(t, email, data.Email)
}

func TestInvalidAccessToken(t *testing.T) {
	jwtService := NewJWT("test-secret")

	isValid, data := jwtService.ParseToken("invalid.token.format")

	require.False(t, isValid)
	require.Nil(t, data)
}

func TestRefreshTokenType(t *testing.T) {
	// Проверяем, что access token нельзя использовать как refresh token
	const email = "test@example.com"
	jwtService := NewJWT("test-secret")

	tokenPair, err := jwtService.GenerateTokenPair(JWTData{
		Email: email,
	})

	require.NoError(t, err)

	// Пытаемся использовать access token как refresh token
	isValid, data := jwtService.ParseRefreshToken(tokenPair.AccessToken)

	if isValid {
		t.Fatal("Access token should not be valid as refresh token")
	}

	if data != nil {
		t.Fatal("Data should be nil when using access token as refresh token")
	}
}

func TestTokenExpiration(t *testing.T) {

	const email = "test@example.com"
	jwtService := NewJWT("test-secret")

	// Устанавливаем время жизни достаточное для выполнения теста
	jwtService.AccessTokenTTL = time.Second

	tokenPair, err := jwtService.GenerateTokenPair(JWTData{
		Email: email,
	})

	require.NoError(t, err)

	// Проверяем, что токен сначала валиден
	isValid, _ := jwtService.ParseToken(tokenPair.AccessToken)
	if !isValid {
		t.Fatal("Token should be valid initially")
	}

	// Ждем истечения срока действия
	time.Sleep(time.Second * 2)

	// Проверяем, что токен стал невалидным
	isValid, _ = jwtService.ParseToken(tokenPair.AccessToken)
	if isValid {
		t.Fatal("Token should be invalid after expiration")
	}
}
