package connection_strings

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/AntoniusIvan/go-Sample3FrlConnect/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

func (db DBClient) getConnectionString() (string, string) {
	var dataSourceName string
	var driverName string

	switch db.Type {
	case "mysql":
		if db.Port == nil {
			temp := 3306
			db.Port = &temp
		}
		if db.Type == "" {
			db.Type = "mysql"
		}
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db.Username, db.Password, db.Host, *db.Port, db.Name)
		driverName = "mysql"
	case "sqlserver":
		fallthrough
	default:
		if db.Port == nil {
			temp := 1433
			db.Port = &temp
		}
		if db.Type == "" {
			db.Type = "sqlserver"
		}

		dataSourceName = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", db.Host, db.Username, db.Password, *db.Port, db.Name)
		driverName = "sqlserver"
	}

	return dataSourceName, driverName
}

func (db DBClient) getDBType() string {
	if strings.Contains(db.ConnString, "MySQL") {
		return "mysql"
	} else if strings.Contains(db.ConnString, "SQLOLEDB") {
		return "sqlserver"
	}

	return "sqlserver"
}

func (db DBClient) ConnectToProjDB() *sqlx.DB { // open db connection
	dataSourceName, driverName := db.getConnectionString()

	logger.Info("Connecting to database: " + dataSourceName)

	client, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = client.PingContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected to database: " + dataSourceName)

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client
}
