package core

import (
	"gorm.io/gorm"
)

var (
	ErrNilDB = gorm.ErrInvalidDB
)

func TableInsert(db *gorm.DB, table string, obj map[string]any) error {
	if db == nil {
		return ErrNilDB
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Table(table).Create(&obj).Error
	})
}

func TableSelect(db *gorm.DB, query string, values ...any) ([]map[string]any, error) {
	if db == nil {
		return nil, ErrNilDB
	}

	data := make([]map[string]any, 0)

	if err := db.Raw(query, values...).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
