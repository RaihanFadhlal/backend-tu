package models

type ProductSafari struct {
	Id                   uint                   `gorm:"primaryKey;autoIncrement"`
	Code                 string                 `gorm:"type:varchar(10);unique;not null"`
	GroupCode            string                 `gorm:"type:varchar(8);not null"`
	Name                 string                 `gorm:"type:varchar(64);not null"`
	Description          string                 `gorm:"type:varchar(256)"`
	Price                int                    `gorm:"type:int;not null"`
	Countries            string                 `gorm:"type:varchar(512);"`
	DayMin               int                    `gorm:"type:int;not null"`
	DayMax               int                    `gorm:"type:int;not null"`
	Contribution         string                 `gorm:"type:product_type;not null;default:'Basic'"`
	Terms                string                 `gorm:"type:varchar(1024)"`
	Image                string                 `gorm:"type:varchar(10)"`
	ProductBenefitSafari []ProductBenefitSafari `gorm:"foreignKey:Code;references:Code"`
	Transaction          []Transaction          `gorm:"foreignKey:ProductCode;references:Code"`
	EnrollmentSafari     []EnrollmentSafari     `gorm:"foreignKey:ProductCode;references:Code"`
	ClaimSafari          []ClaimSafari          `gorm:"foreignKey:ProductCode;references:Code"`
}

type ProductBenefitSafari struct {
	Id          uint   `gorm:"primaryKey;autoIncrement"`
	Code        string `gorm:"type:varchar(10);not null"`
	GroupCode   string `gorm:"type:varchar(10);not null"`
	Description string `gorm:"type:varchar(64);not null"`
	Detail      string `gorm:"type:varchar(64)"`
	Basic       string `gorm:"type:varchar(64)"`
	Gold        string `gorm:"type:varchar(64)"`
	Platinum    string `gorm:"type:varchar(64)"`
	Titanium    string `gorm:"type:varchar(64)"`
}
