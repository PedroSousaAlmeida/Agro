package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	// Carrega .env se existir (ignora erro se não existir)
	godotenv.Load()
}

// Env contém as variáveis de ambiente
type Env struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	// Keycloak
	KeycloakURL      string
	KeycloakRealm    string
	KeycloakClientID string
}

// NewEnv carrega as variáveis de ambiente
func NewEnv() *Env {
	return &Env{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "agro_monitoring"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		// Keycloak
		KeycloakURL:      getEnv("KEYCLOAK_URL", "http://localhost:9090"),
		KeycloakRealm:    getEnv("KEYCLOAK_REALM", "agro-realm"),
		KeycloakClientID: getEnv("KEYCLOAK_CLIENT_ID", "agro-api"),
	}
}

// DSN retorna a connection string do PostgreSQL
func (e *Env) DSN() string {
	return "host=" + e.DBHost +
		" port=" + e.DBPort +
		" user=" + e.DBUser +
		" password=" + e.DBPassword +
		" dbname=" + e.DBName +
		" sslmode=" + e.DBSSLMode
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

// RedisAddr retorna o endereço do Redis
func (e *Env) RedisAddr() string {
	return e.RedisHost + ":" + e.RedisPort
}
