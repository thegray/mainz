package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DBImpl struct {
	config string //todo: ganti config
}

func (d *DBImpl) TestConnect() error {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		"root", "root", "127.0.0.1", "3306", "pbk_main")
	con, err := sql.Open("mysql", connectionString)
	defer con.Close()
	if err != nil {
		log.Printf("[Database] Failed to open DB connection. %+v \n", err)
		return err
	}
	err = con.Ping()
	if err != nil {
		log.Printf("[Database] DB didn't respond. %+v \n", err)
		return err
	}
	return nil
}

func (d *DBImpl) Connect() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		"root", "root", "127.0.0.1", "3306", "pbk_main")
	con, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Printf("[InitializeDatabase] Failed to open DB connection. %+v \n", err)
		return nil, err
	}
	// DB.SetMaxIdleConns(mysqlMaxIdle)
	// DB.SetMaxOpenConns(mysqlMaxOpen)
	// DB.SetConnMaxLifetime(time.Duration(mysqlConnMaxLifetime) * time.Minute)
	return con, err
}
