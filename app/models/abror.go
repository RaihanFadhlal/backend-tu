package models

type ProductAbror struct {
	Id                  uint                  `gorm:"primaryKey;autoIncrement"`
	Code                string                `gorm:"type:varchar(10);unique;not null"`
	GroupCode           string                `gorm:"type:varchar(8);not null"`
	Name                string                `gorm:"type:varchar(64);not null"`
	Description         string                `gorm:"type:varchar(256)"`
	ClaimProcedure      string                `gorm:"type:varchar(1024);"`
	ClosureProcedure    string                `gorm:"type:varchar(1024);"`
	Image               string                `gorm:"type:varchar(10)"`
	Percentage          float32               `gorm:"type:real;not null"`
	RegionCode          string                `gorm:"type:varchar(10);not null"`
	TypeName            string                `gorm:"type:varchar(12);not null"`
	VehicleCode         string                `gorm:"type:varchar(10);not null"`
	AllowedVehicle      string                `gorm:"type:varchar(512);"`
	ProductBenefitAbror []ProductBenefitAbror `gorm:"foreignKey:Code;references:Code"`
	TransactionAbror    []TransactionAbror    `gorm:"foreignKey:ProductCode;references:Code"`
	EnrollmentAbror     []EnrollmentAbror     `gorm:"foreignKey:ProductCode;references:Code"`
	ClaimAbror          []ClaimAbror          `gorm:"foreignKey:ProductCode;references:Code"`
}

type ProductBenefitAbror struct {
	Id          uint   `gorm:"primaryKey;autoIncrement"`
	Code        string `gorm:"type:varchar(10);not null"`
	GroupCode   string `gorm:"type:varchar(8);not null"`
	Type        string `gorm:"type:abror_tcs;not null;default:'Manfaat'"`
	Description string `gorm:"type:varchar(512);not null"`
	Standard    string `gorm:"type:varchar(512)"`
	Premium     string `gorm:"type:varchar(512)"`
}

type TypeAbror struct {
	ID           uint           `gorm:"primaryKey"`
	Type         string         `gorm:"type:varchar(12);unique;not null"`
	BasePrice    int            `gorm:"type:int;not null"`
	ProductAbror []ProductAbror `gorm:"foreignKey:TypeName;references:Type"`
}

type VehicleType struct {
	ID           uint           `gorm:"primaryKey"`
	Code         string         `gorm:"type:varchar(10);unique;not null"`
	Min          int64          `gorm:"type:int;not null"`
	Max          int64          `gorm:"type:int;not null"`
	ProductAbror []ProductAbror `gorm:"foreignKey:VehicleCode;references:Code"`
}

type Region struct {
	ID           uint           `gorm:"primaryKey"`
	Code         string         `gorm:"type:varchar(10);unique;not null"`
	Description  string         `gorm:"type:text"`
	Provinces    string         `gorm:"type:text"`
	Plat         string         `gorm:"type:text"`
	ProductAbror []ProductAbror `gorm:"foreignKey:RegionCode;references:Code"`
}

type Car struct {
	ID              uint              `gorm:"primaryKey"`
	Name            string            `gorm:"type:varchar(16);unique;not null"`
	Brand           string            `gorm:"type:varchar(16);not null"`
	Price           int               `gorm:"type:int;not null"`
	EnrollmentAbror []EnrollmentAbror `gorm:"foreignKey:CarType;references:Name"`
}
