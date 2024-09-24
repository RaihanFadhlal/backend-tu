package models

import (
	"time"

	"github.com/google/uuid"
)

type ClaimSafari struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	ClaimId      string    `gorm:"type:varchar(32);unique;not null"`
	RegistrantId string    `gorm:"type:varchar(50);not null"`
	EnrollmentId string    `gorm:"type:varchar(32);not null"`
	PolicyId     string    `gorm:"type:varchar;not null"`
	ProductCode  string    `gorm:"type:varchar(8);default:'';not null"`
	ProductName  string    `gorm:"type:varchar(64);default:'';not null"`
	Status       string    `gorm:"type:claim_status;default:'Diproses';not null"`
	DateReport   string    `gorm:"type:varchar(10);default:''"`
	DateAccident string    `gorm:"type:varchar(10);default:''"`
	Location     string    `gorm:"type:varchar(32);default:''"`
	Detail       string    `gorm:"type:varchar(1024);default:''"`
	Evidence     string    `gorm:"type:varchar(64);default:''"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp;not null"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp;not null"`
	PayProof     string    `gorm:"type:varchar(64);default:''"`
	CoverCost    int       `gorm:"type:int"`
	Message      string    `gorm:"type:varchar(1024)"`
}

type ClaimAbror struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	ClaimId      string    `gorm:"type:varchar(32);unique;not null"`
	RegistrantId string    `gorm:"type:varchar(50);not null"`
	EnrollmentId string    `gorm:"type:varchar(32);not null"`
	PolicyId     string    `gorm:"type:varchar;not null"`
	ProductCode  string    `gorm:"type:varchar(8);default:'';not null"`
	ProductName  string    `gorm:"type:varchar(64);default:'';not null"`
	Status       string    `gorm:"type:claim_status;default:'Diproses';not null"`
	DateReport   string    `gorm:"type:varchar(10);default:''"`
	DateAccident string    `gorm:"type:varchar(10);default:''"`
	Location     string    `gorm:"type:varchar(32);default:''"`
	Case         string    `gorm:"type:varchar(128);default:''"`
	Detail       string    `gorm:"type:varchar(1024);default:''"`
	Evidence     string    `gorm:"type:varchar(64);default:''"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp;not null"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp;not null"`
	PayProof     string    `gorm:"type:varchar(64);default:''"`
	CoverCost    int       `gorm:"type:int"`
	Message      string    `gorm:"type:varchar(1024)"`
}
