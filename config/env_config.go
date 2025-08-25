package config

import (
	"os"
)

type EnvConfig struct {
	Postgres struct {
		HOST     string
		Database string
		Username string
		Password string
		Port     string
	}

	CORS struct {
		AllowDomains string
		GlobalDomain string
	}
}

func LoadEnvConfig() *EnvConfig {
	var config EnvConfig

	// Postgres
	config.Postgres.HOST = os.Getenv("PGPOOL_HOST")
	config.Postgres.Database = os.Getenv("PGPOOL_DB")
	config.Postgres.Username = os.Getenv("PGPOOL_USER")
	config.Postgres.Password = os.Getenv("PGPOOL_PASSWORD")
	config.Postgres.Port = os.Getenv("PGPOOL_PORT")

	//// JWT
	//config.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")
	//config.JWT.Algorithm = os.Getenv("JWT_ALGORITHM")
	//
	//if val := os.Getenv("JWT_EXPIRE"); val != "" {
	//	fmt.Sscanf(val, "%d", &config.JWT.Expire)
	//} else {
	//	config.JWT.Expire = 3600 * 24 * 7
	//}
	//
	config.CORS.AllowDomains = os.Getenv("ALLOWED_DOMAINS")
	config.CORS.GlobalDomain = os.Getenv("GLOBAL_DOMAIN")

	//config.Redis.Address = os.Getenv("REDIS_ADDRESS")
	//config.Redis.Password = os.Getenv("REDIS_PASSWORD")
	//config.Redis.Database, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
	//if config.Redis.Database == 0 {
	//	config.Redis.Database = 0
	//}

	//config.PrivateKey = os.Getenv("PRIVATE_KEY")
	//
	//config.ExternalService.AuthorizationServiceURL = os.Getenv("AUTHORIZATION_SERVICE_URL")
	//if config.ExternalService.AuthorizationServiceURL == "" {
	//	config.ExternalService.AuthorizationServiceURL = "http://localhost:8080"
	//}
	//config.ExternalService.UploadServiceURL = os.Getenv("UPLOAD_SERVICE_URL")
	//if config.ExternalService.UploadServiceURL == "" {
	//	config.ExternalService.UploadServiceURL = "http://localhost:8081"
	//}
	//config.ExternalService.CDNServiceURL = os.Getenv("CDN_SERVICE_URL")
	//if config.ExternalService.CDNServiceURL == "" {
	//	config.ExternalService.CDNServiceURL = "http://localhost:8082"
	//}

	return &config
}
