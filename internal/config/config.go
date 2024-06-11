package config

import (
	"crypto/rsa"
	"flag"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type AppConfig struct {
	ServerAddress        string
	DatabaseDSN          string
	AccrualSystemAddress string
	PrivateKey           *rsa.PrivateKey
	PublicKey            *rsa.PublicKey
}

var Config AppConfig

func Init() {
	flag.StringVar(&Config.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Config.DatabaseDSN, "d", "user=postgres dbname=postgres sslmode=disable password=postgres", "database dsn")
	flag.StringVar(&Config.AccrualSystemAddress, "r", "http://localhost:8081", "accrual system address")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		Config.ServerAddress = envRunAddr
	}
	if databaseDSN := os.Getenv("DATABASE_URI"); databaseDSN != "" {
		Config.DatabaseDSN = databaseDSN
	}
	if accrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accrualSystemAddress != "" {
		Config.AccrualSystemAddress = accrualSystemAddress
	}
	pwd, _ := os.Getwd()
	pubPEMData, _ := os.ReadFile(pwd + `/internal/config/jwt/public.pem`)
	Config.PublicKey, _ = jwt.ParseRSAPublicKeyFromPEM(pubPEMData)
	privatePEMData, _ := os.ReadFile(pwd + `/internal/config/jwt/private.pem`)
	Config.PrivateKey, _ = jwt.ParseRSAPrivateKeyFromPEM(privatePEMData)
}
