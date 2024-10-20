package main

import (
    "fmt"
    "github.com/go-pg/pg/v10"
    "github.com/brianvoe/gofakeit/v6"
)

// account table model
type Account struct {
    AccID    int64  `pg:"acc_id,pk"`
    Username string `pg:"username"`
    Email    string `pg:"email"`
}

// character table model
type Character struct {
    CharID  int64  `pg:"char_id,pk"`
    AccID   int64  `pg:"acc_id"`
    ClassID int    `pg:"class_id"`
}

// scores table model
type Score struct {
    ScoreID     int64 `pg:"score_id,pk"`
    CharID      int64 `pg:"char_id"`
    RewardScore int   `pg:"reward_score"`
}

func connect() *pg.DB {
    
    db := pg.Connect(&pg.Options{
        User:     "postgres",       
        Password: "oredayak1",  
        Addr:     "localhost:5432", 
        Database: "wira_db", 
    })







///////////////connection test

	var n int
    _, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
    if err != nil {
        fmt.Println("Could not connect to the database:", err)
        return nil
    }
    fmt.Println("Connected to PostgreSQL successfully!")
    return db
}




///////////////creating tables

func createTables(db *pg.DB) {
    // Create Account table
    err := db.Model((*Account)(nil)).CreateTable(&pg.CreateTableOptions{
        IfNotExists: true,
    })
    if err != nil {
        panic(err)
    }

    // Create Character table
    err = db.Model((*Character)(nil)).CreateTable(&pg.CreateTableOptions{
        IfNotExists: true,
    })
    if err != nil {
        panic(err)
    }

    // Create Scores table
    err = db.Model((*Score)(nil)).CreateTable(&pg.CreateTableOptions{
        IfNotExists: true,
    })
    if err != nil {
        panic(err)
    }

    fmt.Println("Tables created successfully!")
}





/////////////////gofakeit

func generateFakeData(db *pg.DB) {

    for i := 0; i < 100000; i++ {
        account := &Account{
            Username: gofakeit.Username(),
            Email:    gofakeit.Email(),
        }
        _, err := db.Model(account).Insert()
        if err != nil {
            fmt.Println("Error inserting account:", err)
            return
        }

        //each account, generate 8 characters(classes)

        for j := 0; j < 8; j++ {
            character := &Character{
                AccID:   account.AccID,
                ClassID: gofakeit.Number(1, 8), // Random class ID between 1 and 8
            }
            _, err := db.Model(character).Insert()
            if err != nil {
                fmt.Println("Error inserting character:", err)
                return
            }

            //each character, generate random scores

            for k := 0; k < 10; k++ {
                score := &Score{
                    CharID:      character.CharID,
                    RewardScore: gofakeit.Number(100, 1000), // Random score between 100 and 1000
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
