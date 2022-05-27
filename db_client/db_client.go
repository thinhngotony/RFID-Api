package db_client

import (
	"context"
	"database/sql"
	f "fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "Rfid@2022"
	hostname = "192.168.1.244:3306"
	dbname   = "RFID"
)

func dsn(dbName string) string {
	return f.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func DbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return nil, err
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return nil, err
	}
	log.Printf("rows affected %d\n", no)

	db.Close()
	db, err = sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening database", err)
		return nil, err
	}

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging database", err)
		return nil, err
	}
	log.Printf("Verified connection from %s database with Ping\n", dbname)
	return db, nil
}

func InsertToTable(db *sql.DB) error {
	insert, err := db.Query("INSERT INTO `RFID`.`Covert_RFID_JANCODE` (`drgm_rfid_cd`, `drgm_jan`, `drgm_jan2`) VALUES ('Test', 'Test', 'Test');")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	f.Println("Successful Insert to Database!")
	return nil
}

func ConvertFromRFID(db *sql.DB, rfid_code string) (string, string, error) {
	log.Printf("Getting JAN code")
	query := `select drgm_jan, drgm_jan2 from Covert_RFID_JANCODE where drgm_rfid_cd = ?;`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return "", "", err
	}
	defer stmt.Close()
	var jan_code_1, jan_code_2 string
	row := stmt.QueryRowContext(ctx, rfid_code)
	if err := row.Scan(&jan_code_1, &jan_code_2); err != nil {
		return "", "", err
	}
	return jan_code_1, jan_code_2, nil

}
