package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Config struct {
	DBDriver   string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
}

func init() {
	conf := Config{
		DBDriver:   "mysql",
		DBHost:     "the-db.c5balohade6x.ap-southeast-1.rds.amazonaws.com",
		DBUser:     "admin",
		DBPassword: "taptaptap",
		DBName:     "estl",
	}
	dbInfo := fmt.Sprintf("%s:%s@(%s)/%s?loc=UTC&charset=utf8mb4,utf8&parseTime=True", conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBName)

	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open a DB conn")
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if db != nil {
		boil.SetDB(db)
		boil.DebugMode = false
	}
}
