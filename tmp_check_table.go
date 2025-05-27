package main

import (
	"fmt"
	"go-api/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config.InitConfig()
	cfg := config.Get()
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	var columns []struct {
		ColumnName string `gorm:"column:column_name"`
		DataType   string `gorm:"column:data_type"`
	}

	db.Raw("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'users' ORDER BY ordinal_position").Scan(&columns)

	fmt.Println("Users table structure:")
	for _, col := range columns {
		fmt.Printf("  %s: %s\n", col.ColumnName, col.DataType)
	}
}
