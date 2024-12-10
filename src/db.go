package main

import (
    "os"
    "log"
    "database/sql"
    "github.com/go-sql-driver/mysql"
)

var DB *sql.DB;

func init() {
    cfg := mysql.Config {
        User: os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net: "tcp",
        Addr: "127.0.0.1:3306",
        DBName: "electro",
    };

    var err error;
    DB, err = sql.Open("mysql", cfg.FormatDSN());

    if err != nil {
        log.Fatalf("Couldn't connect to db: %s", err.Error());
    }
    
    err = DB.Ping();

    if err != nil {
        log.Fatalf("Db ping failed: %s", err.Error());
    }

    log.Println("Db connected");
}
