package controllers

import (
	"main/db_client"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Post struct {
	RFID     string `json: "rfid`
	JanCode1 string `json: "jancode_1`
	JanCode2 string `json: "jancode_2`
}

func CreatePost(c *gin.Context) {
	var reqBody Post
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	db, _ := db_client.DbConnection()
	jancode_1, jan_code_2, _ := db_client.ConvertFromRFID(db, reqBody.RFID)

	c.JSON(http.StatusCreated, gin.H{
		"error":     false,
		"rfid":      reqBody.RFID,
		"jancode_1": jancode_1,
		"jancode_2": jan_code_2,
	})

}
