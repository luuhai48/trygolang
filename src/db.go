package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	gm "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {
	uri := GetEnv("DB_CONNECTION")
	db, err := gorm.Open(gm.Open(uri), &gorm.Config{
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

func MakeMigration(name string) {
	seqDigits := 6
	ext := ".sql"
	dir := filepath.Clean("migrations")

	var version string

	files, err := filepath.Glob(filepath.Join(dir, "*"+ext))
	if err != nil {
		panic(err)
	}
	version, err = nextSeqVersion(files, seqDigits)
	if err != nil {
		panic(err)
	}

	versionGlob := filepath.Join(dir, version+"_*"+ext)
	matches, err := filepath.Glob(versionGlob)
	if err != nil {
		panic(err)
	}

	if len(matches) > 0 {
		panic("Duplicate migration version: " + version)
	}
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s%s", version, name, direction, ext)
		filename := filepath.Join(dir, basename)

		if err = createFile(filename); err != nil {
			panic(err)
		}

		absPath, _ := filepath.Abs(filename)
		log.Println(absPath)
	}
}

func nextSeqVersion(matches []string, seqDigits int) (string, error) {
	if seqDigits <= 0 {
		return "", errors.New("digits must be positive")
	}

	nextSeq := uint64(1)

	if len(matches) > 0 {
		filename := matches[len(matches)-1]
		matchSeqStr := filepath.Base(filename)
		idx := strings.Index(matchSeqStr, "_")

		if idx < 1 { // Using 1 instead of 0 since there should be at least 1 digit
			return "", fmt.Errorf("malformed migration filename: %s", filename)
		}

		var err error
		matchSeqStr = matchSeqStr[0:idx]
		nextSeq, err = strconv.ParseUint(matchSeqStr, 10, 64)

		if err != nil {
			return "", err
		}

		nextSeq++
	}

	version := fmt.Sprintf("%0[2]*[1]d", nextSeq, seqDigits)

	if len(version) > seqDigits {
		return "", fmt.Errorf("next sequence number %s too large. At most %d digits are allowed", version, seqDigits)
	}

	return version, nil
}

func createFile(filename string) error {
	// create exclusive (fails if file already exists)
	// os.Create() specifies 0666 as the FileMode, so we're doing the same
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {
		return err
	}

	return f.Close()
}

func RunMigration() {
	log.Println("Running migration...")

	uri := GetEnv("DB_CONNECTION") + "&multiStatements=true"
	db, err := sql.Open("mysql", uri)
	if err != nil {
		panic(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		panic(err)
	}

	m.Up()
	db.Close()
}
