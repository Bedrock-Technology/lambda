package model

import (
	"gorm.io/gorm"
)

type ClaimInfo struct {
	gorm.Model

	Address        string `gorm:"type:varchar(255);not null;uniqueIndex"`
	CexType        string `gorm:"type:varchar(255);not null;index"`
	CexUid         string `gorm:"type:varchar(255);not null"`
	DepositAddress string `gorm:"type:varchar(255);not null"`
	Signature      string `gorm:"type:varchar(255);not null"`
}

func (ClaimInfo) TableName() string {
	return "claim_info"
}
