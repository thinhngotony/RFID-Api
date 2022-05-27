package main

import (
	"database/sql"
	"log"
	"main/controllers"
	"main/db_client"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db_client.DbConnection()
	if err != nil {
		log.Printf("Error %s when getting db connection", err)
		return
	}
	defer db.Close()
	log.Printf("Successfully connected to database")

	rfid_code := "0xRFID"
	jan_code_1, jan_code_2, err := db_client.ConvertFromRFID(db, rfid_code)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("Product %s not found in DB", rfid_code)
	case err != nil:
		log.Printf("Encountered err %s when fetching price from DB", err)
	default:
		log.Printf("JanCode 1 of %s is %s, JanCode 2 is %s", rfid_code, jan_code_1, jan_code_2)
	}

	r := gin.Default()
	r.POST("/", controllers.CreatePost)
	if err := r.Run(":333"); err != nil {
		panic(err.Error())
	}

}
