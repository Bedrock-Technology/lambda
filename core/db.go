package core

import (
	"gorm.io/gorm"
)

func TableInsert(db *gorm.DB, table string, obj map[string]any) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Table(table).Create(&obj).Error
	})
}

func TableSelect(db *gorm.DB, query string) ([]map[string]any, error) {
	data := make([]map[string]any, 0)

	if err := db.Raw(query).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
