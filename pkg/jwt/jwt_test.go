package jwt

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// Загружаем переменные из .env файла
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
}

func TestGenerateTokenPair(t *testing.T) {
	const email = "test@example.com"
	jwtService := NewJWT(os.Getenv("KEY"))

	// Устанавливаем короткое время жизни для тестов
	jwtService.AccessTokenTTL = time.Second * 10
	jwtService.RefreshTokenTTL = time.Second * 20

	tokenPair, err := jwtService.GenerateTokenPair(JWTData{
		Email: email,
	})

	if err != nil {
		t.Fatal(err)
	}

	if tokenPair.AccessToken == "" {
		t.Fatal("Access token is empty")
	}

	if tokenPair.RefreshToken == "" {
		t.Fatal("Refresh token is empty")
	}
}

func TestParseAccessToken(t *testing.T) {
	const email = "test@example.com"
	jwtService := NewJWT(os.Getenv("KEY"))

	tokenPair, err := jwtService.GenerateTokenPair(JWTData{
		Email: email,
	})

	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.ParseToken(tokenPair.AccessToken)

	if !isValid {
		t.Fatal("Access token is not valid")
	}

	if data.Email != email {
		t.Fatalf("Email %s != %s", data.Email, email)
	}
}

func TestParseRefreshToken(t *testing.T) {
	const email = "test@example.com"
	jwtService := NewJWT(os.Getenv("KEY"))

	tokenPair, err := jwtService.GenerateTokenPair(JWTData{
		Email: email,
	})

	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.ParseRefreshToken(tokenPair.RefreshToken)

	if !isValid {
		t.Fatal("Refresh token is not valid")
	}

	if data.Email != email {
		t.Fatalf("Email %s != %s", data.Email, email)
	}
}

func TestInvalidAccessToken(t *testing.T) {
	jwtService := NewJWT(os.Getenv("KEY"))

	isValid, data := jwtService.ParseToken("invalid.token.format")

	if isValid {
		t.Fatal("Invalid token should not be valid")
	}

	if data != nil {
		t.Fatal("Data should be nil for invalid token")
	}
}

func TestRefreshTokenType(t *testing.T) {
	// Проверяем, что access token нельзя использовать как refresh token
	const email = "test@example.com"
	jwtService := NewJWT(os.Getenv("KEY"))

	tokenPair, err := jwtService.GenerateTokenPair(JWTData{
		Email: email,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Пытаемся использовать access token как refresh token
	isValid, data := jwtService.ParseRefreshToken(tokenPair.AccessToken)

	if isValid {
		t.Fatal("Access token should not be valid as refresh token")
	}

	if data != nil {
		t.Fatal("Data should be nil when using access token as refresh token")
	}
}
