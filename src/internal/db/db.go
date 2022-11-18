package db

import (
	"database/sql"
	"fmt"
	"goduckduckgo/internal/models"
	"goduckduckgo/pkg/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	gormClient *gorm.DB
	sqlClient  *sql.DB
}

// Returns a pointer to the newly created Server
// GORM allows to initialize *gorm.DB with an existing database connection
func NewDB(cfg config.Config) *Database {

	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.StoreConfig.DBHost, cfg.StoreConfig.DBUser, cfg.StoreConfig.DBPassword, cfg.StoreConfig.DBName, cfg.StoreConfig.DBPort)

	//The Open function should be called just once. It is rarely necessary to close a DB.
	sqlDB, err := sql.Open("pgx", DSN)

	if err != nil {
		log.Fatal("Cannot connect to DB")
	}

	//Initialize db session based on dialector
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:  DSN,
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Error occurred while connecting with the database")
	}

	//Configure DB params
	//GORM using database/sql to maintain connection pool
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("DB INIT!!")
	return &Database{
		gormClient: gormDB,
		sqlClient:  sqlDB,
	}
}

func (d *Database) GetDatabaseConnection() error {

	if err := d.sqlClient.Ping(); err != nil {
		return err
	}

	return nil
}

func (d *Database) CloseDBConnection() error {

	if err := d.sqlClient.Close(); err != nil {
		return err
	}

	return nil
}

func (d *Database) AutoMigrateDB() error {

	if err := d.gormClient.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	return nil
}
