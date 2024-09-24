package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID               uuid.UUID          `gorm:"type:uuid;primaryKey"`
	TransactionId    string             `gorm:"type:varchar(16);unique;not null"`
	RegistrantId     string             `gorm:"type:varchar(50);not null"`
	ProductCode      string             `gorm:"type:varchar(8);not null"`
	ProductName      string             `gorm:"type:varchar(64);not null"`
	ProductPrice     int                `gorm:"type:int;not null"`
	Capacity         int                `gorm:"type:int;default:1;not null"`
	TotalPrice       int                `gorm:"type:int;not null"`
	Status           string             `gorm:"type:status_type;default:'Menunggu Pembayaran'"`
	CreatedAt        time.Time          `gorm:"type:timestamp;default:current_timestamp;not null"`
	ExpiredAt        time.Time          `gorm:"type:timestamp;default:current_timestamp;not null"`
	EnrollmentSafari []EnrollmentSafari `gorm:"foreignKey:TransactionId;references:TransactionId"`
}

type EnrollmentSafari struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey"`
	EnrollmentId  string        `gorm:"type:varchar(16);unique;not null"`
	RegistrantId  string        `gorm:"type:varchar(50);not null"`
	TransactionId string        `gorm:"type:varchar(14);not null"`
	CreatedAt     time.Time     `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt     time.Time     `gorm:"type:timestamp;default:current_timestamp"`
	Phone         string        `gorm:"type:varchar(14);default:''"`
	ProductCode   string        `gorm:"type:varchar(8);default:''"`
	ProductName   string        `gorm:"type:varchar(64);default:''"`
	From          string        `gorm:"type:varchar(32);default:''"`
	Destination   string        `gorm:"type:varchar(32);default:''"`
	DateStart     string        `gorm:"type:varchar(10);default:''"`
	DateEnd       string        `gorm:"type:varchar(10);default:''"`
	Contribution  string        `gorm:"type:varchar(10);default:''"`
	Capacity      int           `gorm:"type:int;default:1"`
	Name          string        `gorm:"type:varchar(64);not null"`
	Birthdate     string        `gorm:"type:varchar(10);default:''"`
	Birthplace    string        `gorm:"type:varchar(20);default:''"`
	Gender        string        `gorm:"type:user_gender;default:'M'"`
	Passport      string        `gorm:"type:varchar(12);default:''"`
	PolicyId      string        `gorm:"type:varchar(16);default:''"`
	ClaimSafari   []ClaimSafari `gorm:"foreignKey:EnrollmentId;references:EnrollmentId"`
}

type TransactionAbror struct {
	ID              uuid.UUID         `gorm:"type:uuid;primaryKey"`
	TransactionId   string            `gorm:"type:varchar(16);unique;not null"`
	RegistrantId    string            `gorm:"type:varchar(50);not null"`
	ProductCode     string            `gorm:"type:varchar(8);not null"`
	ProductName     string            `gorm:"type:varchar(64);not null"`
	ProductPrice    int               `gorm:"type:int;not null"`
	Capacity        int               `gorm:"type:int;default:1;not null"`
	TotalPrice      int               `gorm:"type:int;not null"`
	Status          string            `gorm:"type:status_type;default:'Menunggu Pembayaran'"`
	CreatedAt       time.Time         `gorm:"type:timestamp;default:current_timestamp;not null"`
	ExpiredAt       time.Time         `gorm:"type:timestamp;default:current_timestamp;not null"`
	EnrollmentAbror []EnrollmentAbror `gorm:"foreignKey:TransactionId;references:TransactionId"`
}

type EnrollmentAbror struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey"`
	EnrollmentId  string       `gorm:"type:varchar(16);unique;not null"`
	RegistrantId  string       `gorm:"type:varchar(50);not null"`
	TransactionId string       `gorm:"type:varchar(14);not null"`
	CreatedAt     time.Time    `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt     time.Time    `gorm:"type:timestamp;default:current_timestamp"`
	Phone         string       `gorm:"type:varchar(14);default:'';not null"`
	ProductCode   string       `gorm:"type:varchar(8);default:'';not null"`
	ProductName   string       `gorm:"type:varchar(64);default:'';not null"`
	CarBrand      string       `gorm:"type:varchar(16);not null"`
	CarType       string       `gorm:"type:varchar(16);not null"`
	Year          string       `gorm:"type:varchar(4);not null"`
	Plat          string       `gorm:"type:varchar(10);not null"`
	Chassis       string       `gorm:"type:varchar(17);not null"`
	Engine        string       `gorm:"type:varchar(18);not null"`
	Image1        string       `gorm:"type:varchar(64);default:'';not null"`
	Image2        string       `gorm:"type:varchar(64);default:'';not null"`
	Image3        string       `gorm:"type:varchar(64);default:'';not null"`
	Image4        string       `gorm:"type:varchar(64);default:'';not null"`
	DateStart     string       `gorm:"type:varchar(10);default:'';not null"`
	DateEnd       string       `gorm:"type:varchar(10);default:'';not null"`
	Contribution  string       `gorm:"type:varchar(10);default:'';not null"`
	Name          string       `gorm:"type:varchar(64);not null"`
	Birthdate     string       `gorm:"type:varchar(10);default:'';not null"`
	Birthplace    string       `gorm:"type:varchar(20);default:'';not null"`
	Gender        string       `gorm:"type:user_gender;default:'M'"`
	PolicyId      string       `gorm:"type:varchar(16);default:''"`
	IdUser        string       `gorm:"type:varchar(64);default:'';not null"`
	ClaimAbror    []ClaimAbror `gorm:"foreignKey:EnrollmentId;references:EnrollmentId"`
} 
