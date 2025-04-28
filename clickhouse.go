package main

import (
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	clickhouseDriver "gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

var (
	clickhouseDB = make(map[string]*gorm.DB)
)

func loadClickhouse(key, dsn string) error {
	opt, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		return err
	}

	opt.Settings = clickhouse.Settings{
		"max_execution_time": 60,
	}
	opt.DialTimeout = 15 * time.Second
	opt.Compression = &clickhouse.Compression{
		Method: clickhouse.CompressionLZ4,
	}

	d := clickhouse.OpenDB(opt)
	d.SetMaxIdleConns(10)
	d.SetMaxOpenConns(100)
	d.SetConnMaxLifetime(120 * time.Minute)

	db, err := gorm.Open(clickhouseDriver.New(clickhouseDriver.Config{Conn: d}))
	if err != nil {
		return err
	}

	clickhouseDB[key] = db
	return nil
}
