package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {
	uri := GetEnv("DB_CONNECTION")
	db, err := gorm.Open(mysql.Open(uri), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Panic("Error connecting to database: " + err.Error())
	}

	DB = db
}

func CloseDatabase() {
	if DB != nil {
		log.Println("Disconnecting from database")
		db, _ := DB.DB()
		db.Close()
	}
}

func DropUnusedColumns(dsts ...interface{}) {
	for _, dst := range dsts {
		stmt := &gorm.Statement{DB: DB}
		stmt.Parse(dst)
		fields := stmt.Schema.Fields
		columns, _ := DB.Migrator().ColumnTypes(dst)
		for i := range columns {
			found := false
			for j := range fields {
				if columns[i].Name() == fields[j].DBName {
					found = true
					break
				}
			}
			if !found {
				DB.Debug().Migrator().DropColumn(dst, columns[i].Name())
			}
		}
	}
}
