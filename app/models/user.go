package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;unique"`
	Type              string             `gorm:"type:user_type;default:'user'"`
	Name              string             `gorm:"type:varchar(64);not null"`
	Email             string             `gorm:"type:varchar(50);unique;not null"`
	Password          string             `gorm:"type:varchar;not null"`
	CreatedAt         time.Time          `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt         time.Time          `gorm:"type:timestamp;default:current_timestamp"`
	VerificationToken string             `gorm:"type:varchar(255);default:''"`
	RefreshToken      string             `gorm:"type:varchar(255);default:''"`
	IsVerified        bool               `gorm:"default:false"`
	Birthdate         string             `gorm:"type:varchar(10);default:''"`
	Phone             string             `gorm:"type:varchar(15);default:''"`
	Gender            string             `gorm:"type:user_gender;default:'M'"`
	Birthplace        string             `gorm:"type:varchar(20);default:''"`
	Address           string             `gorm:"type:varchar(64);default:''"`
	Image             string             `gorm:"type:varchar(64);default:''"`
	Transaction       []Transaction      `gorm:"foreignKey:RegistrantId;references:Email"`
	TransactionAbror  []TransactionAbror `gorm:"foreignKey:RegistrantId;references:Email"`
	EnrollmentSafari  []EnrollmentSafari `gorm:"foreignKey:RegistrantId;references:Email"`
	EnrollmentAbror   []EnrollmentAbror  `gorm:"foreignKey:RegistrantId;references:Email"`
	ClaimSafari       []ClaimSafari      `gorm:"foreignKey:RegistrantId;references:Email"`
	ClaimAbror        []ClaimAbror       `gorm:"foreignKey:RegistrantId;references:Email"`
}
