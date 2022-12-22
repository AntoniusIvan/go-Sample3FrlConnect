package main

import (
	"database/sql"
	"fmt"
	"go-Sample3FrlConnect/connection_strings"
	"log"

	"github.com/gin-gonic/gin"
	mssql "github.com/microsoft/go-mssqldb"
)

var db *sql.DB

func main() {
	//connection_strings.ConnectToDB("TTSDB", "TestaPub")
	connection_strings.ConnectToDB("TTSDB", "TestPub")

	test1 := fmt.Sprintf("server=%s;",
		mssql.ErrorTypeSliceIsEmpty)
	connString := test1
	connString = fmt.Sprintf("sqlserver://sa:sqlserver@192.168.137.99:49173?database=OMSDB&connection+timeout=30")
	var err error
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	//COMMENTED 2022-12-12 15:03
	router := gin.Default()

	//router.GET("/guestmodelbooks", routes.GetModelGuestBook)
	//router.POST("/guestbook/modelcreate", routes.AddGuestBook)
	//router.POST("/guestbook/create", routes.CanPurchase(3))

	router.Run("localhost:8005")

}
