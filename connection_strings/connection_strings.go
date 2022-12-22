package connection_strings

//first entry.

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"go-Sample3FrlConnect/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

type DBConfig struct {
	DbConnection string
	DbHost       string
	DbPort       string
	DbDatabase   string
	DbUsername   string
	DbPassword   string
}

type DBClient struct {
	Username string  `json:"usn" db:"usn"`
	Password string  `json:"pwd" db:"pwd"`
	Name     string  `json:"name" db:"ds_cd"`
	Host     string  `json:"dns" db:"dns"`
	Port     *int    `json:"port_no" db:"port_no"`
	Instance *string `json:"instance" db:"instn_nm"`
	Type     string  `json:"type"`
	//ConnString    string  `json:"connection_string"`
	ConnString string `json:"connection_string" db:"attr_go"`
}

func ConnectToDB(projCode, pubKey string) (*DBClient, error) {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	logger.Fatal("Error loading .env file")
	// 	panic(err)
	// }

	mainDB := connectToMainDatabase(pubKey)

	projDB, err := getDatabaseDetailByProjectCode(mainDB, projCode, pubKey)
	if err != nil {
		return nil, err
	}

	//dbClient := projDB.ConnectToDB()

	return projDB, nil
}

func getDatabaseDetailByProjectCode(mainDB *sqlx.DB, projCode, pubKey string) (*DBClient, error) {
	var dbClient DBClient
	err := mainDB.Get(&dbClient, "SELECT usn, pwd, ds_cd, dns, port_no, instn_nm, attr_go FROM vw_DataSource WHERE ds_cd = @p1", projCode)
	if err != nil {
		logger.Error(err.Error())
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no code named: %s", projCode)
		}
		return nil, err
	}

	// TODO get new column of type
	dbClient.Type = dbClient.getDBType()

	//dataSourceName, _ := dbClient.getConnectionString()

	key, err := decryptCryptoKey(pubKey)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	dbClient.Password, err = tripleDESECBDecrypt(dbClient.Password, key)

	//dbClient.ConnString = dataSourceName

	return &dbClient, nil
}

func connectToMainDatabase(pubKey string) *sqlx.DB {
	dbUsername := config["DB_USERNAME"]
	dbPassword := config["DB_PASSWORD"]
	dbName := config["DB_DATABASE"]
	dbHost := config["DB_HOST"]
	dbPort := config["DB_PORT"]
	dbInstance := config["DB_INSTANCE"]
	dbConnection := config["DB_CONNECTION"]

	key, err := decryptCryptoKey(pubKey)
	if err != nil {
		log.Fatal("cannot connect to database: ", err.Error())
	}

	password, err := tripleDESECBDecrypt(dbPassword, key)
	if err != nil {
		log.Fatal("cannot connect to database: ", err.Error())
	}

	// open db connection
	dataSourceName := fmt.Sprintf("server=%s\\%s;user id=%s;password=%s;port=%s;database=%s;encrypt=disable;", dbHost, dbInstance, dbUsername, password, dbPort, dbName)
	logger.Info("Connecting to database: " + dataSourceName)

	client, err := sqlx.Open(dbConnection, dataSourceName)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}

	ctx := context.Background()
	err = client.PingContext(ctx)
	if err != nil {
		log.Fatal("Error ping connection pool: ", err.Error())
	}
	fmt.Printf("Connected to database: " + dataSourceName)

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client
}
