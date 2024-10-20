package main

import (
    "fmt"
    "github.com/go-pg/pg/v10"
    "github.com/brianvoe/gofakeit/v6"
)

///////////////////////////Accounts table model

type Accounts struct {
    AccID    int64  `pg:"acc_id,pk"`
    Username string `pg:"username"`
    Email    string `pg:"email"`
}

///////////////////////////Characters table model

type Characters struct {
    CharID  int64  `pg:"char_id,pk"`
    AccID   int64  `pg:"acc_id"`
    ClassID int    `pg:"class_id"`
}

///////////////////////////Scores table model

type Scores struct {
    ScoreID     int64 `pg:"score_id,pk"`
    CharID      int64 `pg:"char_id"`
    RewardScore int   `pg:"reward_score"`
}

///////////////////////////Connect to PostgreSQL

func connect() *pg.DB {
    db := pg.Connect(&pg.Options{
        User:     "postgres",       
        Password: "oredayak1",      
        Addr:     "localhost:5432", 
        Database: "wira_db",        
    })

    /////////////////////////////////////connection test

    var n int
    _, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
    if err != nil {
        fmt.Println("Could not connect to the database:", err)
        return nil
    }
    fmt.Println("Connected to PostgreSQL successfully!")
    return db
}






///////////////////////////////create tables in the db


func createTables(db *pg.DB) {
	
    //Accounts table

    err := db.Model((*Accounts)(nil)).CreateTable(nil)
    if err != nil {
        panic(err)
    }

    //Characters table

    err = db.Model((*Characters)(nil)).CreateTable(nil)
    if err != nil {
        panic(err)
    }

    //Scores table

    err = db.Model((*Scores)(nil)).CreateTable(nil)
    if err != nil {
        panic(err)
    }

    fmt.Println("Tables created successfully!")
}


//////////////////////////////generate 100,000 data with gofakeit and insert into wira_db


func generateFakeData(db *pg.DB) {
    for i := 0; i < 100000; i++ {
        account := &Accounts{
            Username: gofakeit.Username(),
            Email:    gofakeit.Email(),
        }
        _, err := db.Model(account).Insert()
        if err != nil {
            fmt.Println("Error inserting account:", err)
            return
        }

        accID := account.AccID

        ///////////////////////////each account, generate 8 characters (classes)

        for j := 0; j < 8; j++ {
            character := &Characters{
                AccID:   accID,
                ClassID: gofakeit.Number(1, 8),
            }

            _, err := db.Model(character).Insert()
            if err != nil {
                fmt.Println("Error inserting character:", err)
                return
            }

            ///////////////////////////each character, generate scores 100-1000

            for k := 0; k < 10; k++ {
                score := &Scores{
                    CharID:      character.CharID,
                    RewardScore: gofakeit.Number(100, 1000),
                }
                _, err := db.Model(score).Insert()
                if err != nil {
                    fmt.Println("Error inserting score:", err)
                    return
                }
            }
        }
    }
    fmt.Println("Data generation complete!")
}

func main() {
    db := connect()
    if db == nil {
        return /////////////////// Exit if connection failed
    }

    createTables(db)
    generateFakeData(db)
}
