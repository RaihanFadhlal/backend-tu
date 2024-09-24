package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"backendtku/app/models"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string
	DB      *gorm.DB
	BaseUrl string
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(
		&models.User{}, 
		&models.ProductSafari{}, 
		&models.ProductBenefitSafari{}, 
		&models.Transaction{}, 
		&models.EnrollmentSafari{}, 
		&models.ClaimSafari{},
		&models.TypeAbror{},
		&models.VehicleType{},
		&models.Region{},
		&models.ProductAbror{},
		&models.ProductBenefitAbror{},
		&models.TransactionAbror{},
		&models.Car{},
		&models.EnrollmentAbror{},
		&models.ClaimAbror{}, 
	)

	return &Config{
		AppName: viper.GetString("APP_NAME"),
		AppEnv:  viper.GetString("APP_ENV"),
		AppPort: viper.GetString("APP_PORT"),
		BaseUrl: viper.GetString("BASE_URL"),
		DB:      db,
	}, nil
}
