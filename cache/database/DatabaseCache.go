package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql" //this is to import mysql database driver
	"github.com/magiconair/properties"
)

var pool *sql.DB

// Pool to get the database connection pool initiated on startup
func Pool() (result *sql.DB) {
	result = pool
	return
}

// Init is used to initialize database connection pool to be used in the app lifetime
func Init() {
	var err error
	path, err := filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
	if err != nil {
		log.Fatalf("unable to get root directory: %v", err)
	}
	config := properties.MustLoadFile(path+"/properties/database/database.conf", properties.UTF8)
	user := config.GetString("user", "root")
	password := config.GetString("password", "amper123")
	host := config.GetString("host", "127.0.0.1")
	port := config.GetString("port", "3306")
	dbname := config.GetString("dbname", "amper")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)
	pool, err = sql.Open("mysql", dataSourceName) //"root:@tcp(127.0.0.1:3306)/amper"

	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal("unable to use data source name", err)
	}
	//defer pool.Close()

	pool.SetConnMaxLifetime(config.GetDuration("connMaxLifeTime", time.Hour))
	pool.SetMaxIdleConns(config.GetInt("maxIdleConns", 3))
	pool.SetMaxOpenConns(config.GetInt("maxOpenConns", 3))

	/*ctx, stop := context.WithCancel(context.Background())
	defer stop()

	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		select {
		case <-appSignal:
			stop()
		}
	}()

	ping(ctx)*/
}

// Ping the database to verify DSN provided by the user is valid and the
// server accessible. If the ping fails exit the program with an error.
func ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}
