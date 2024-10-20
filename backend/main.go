package main

import (
    "fmt"
    "math"
    "net/http"
    "strconv"
    "time"

    "github.com/brianvoe/gofakeit/v6"
    "github.com/gin-gonic/gin"
    "github.com/go-pg/pg/v10"
    "github.com/patrickmn/go-cache"
)

/////////////////////////// 'Accounts' table model

type Accounts struct {
    AccID    int64  `pg:"acc_id,pk"`
    Username string `pg:"username"`
    Email    string `pg:"email"`
}

/////////////////////////// 'Characters' table model

type Characters struct {
    CharID  int64  `pg:"char_id,pk"`
    AccID   int64  `pg:"acc_id"`
    ClassID int    `pg:"class_id"`
}

/////////////////////////// 'Scores' table model

type Scores struct {
    ScoreID     int64 `pg:"score_id,pk"`
    CharID      int64 `pg:"char_id"`
    RewardScore int   `pg:"reward_score"`
}

/////////////////////////// connect to PostgreSQL

func connect() *pg.DB {
    db := pg.Connect(&pg.Options{
        User:     "postgres",       
        Password: "oredayak1",      
        Addr:     "localhost:5432", 
        Database: "wira_db",        
    })

    ///////////////////////////////////// connection test

    var n int
    _, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
    if err != nil {
        fmt.Println("Could not connect to the database:", err)
        return nil
    }
    fmt.Println("Connected to PostgreSQL successfully!")
    return db
}

/////////////////////////////// create tables in the db


func createTables(db *pg.DB) {
    // Accounts table
    err := db.Model((*Accounts)(nil)).CreateTable(nil)
    if err != nil {
        panic(err)
    }

    // Characters table
    err = db.Model((*Characters)(nil)).CreateTable(nil)
    if err != nil {
        panic(err)
    }

    // Scores table
    err = db.Model((*Scores)(nil)).CreateTable(nil)
    if err != nil {
        panic(err)
    }

    fmt.Println("Tables created successfully!")
}


///////////////////////////// generate fake data


func generateFakeData(db *pg.DB) {
    //////////////////////////// check if accounts already exist
    count, err := db.Model((*Accounts)(nil)).Count()
    if err != nil {
        fmt.Println("Error counting accounts:", err)
        return
    }
    if count > 0 {
        fmt.Println("Data already exists, skipping data generation.")
        return
    }

    ////////////////////////////// generate 100,000 data with gofakeit and insert into wira_db
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

        /////////////////////////// each account, generate 8 characters (classes)
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

            /////////////////////////// each character, generate scores 100-1000
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



////////////////////////////////////// PAGINATION ///////////////////////////////////////////


func paginatedAccounts(db *pg.DB, c *gin.Context) ([]Accounts, int, int, error) {
    page, _ := strconv.Atoi(c.Query("page"))
    limit, _ := strconv.Atoi(c.Query("limit"))
    search := c.Query("search") // Get the search term

    if page == 0 {
        page = 1
    }
    if limit == 0 {
        limit = 10
    }

    var accounts []Accounts
    offset := (page - 1) * limit

    // Build the query with optional search filtering
    query := db.Model(&accounts).Limit(limit).Offset(offset)

    if search != "" {
        query.Where("username ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
    }

    // Fetch accounts with pagination and search
    err := query.Select()
    if err != nil {
        return nil, 0, 0, err
    }

    // Count total accounts for pagination (considering search)
    total, err := db.Model(&Accounts{}).Where("username ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%").Count()
    if err != nil {
        return nil, 0, 0, err
    }

    totalPages := int(math.Ceil(float64(total) / float64(limit)))
    return accounts, total, totalPages, nil
}



////////////////////////////////////// PAGINATION ///////////////////////////////////////////






func main() {
    db := connect()
    if db == nil {
        return /////////////////// Exit if connection failed
    }

    createTables(db)
    generateFakeData(db)








    ////////////////////////////////////////// CACHING //////////////////////////////////////////

    
    cacheInstance := cache.New(5*time.Minute, 10*time.Minute)

    r := gin.Default()

    ////////////////////////////////// Route for paginated accounts with caching

    r.GET("/accounts", func(c *gin.Context) {

        // Create a cache key based on the page, limit, and search query parameters

        cacheKey := fmt.Sprintf("accounts_page_%s_limit_%s_search_%s", c.Query("page"), c.Query("limit"), c.Query("search"))
    
        // Check if cached data exists

        if cachedData, found := cacheInstance.Get(cacheKey); found {
            c.JSON(http.StatusOK, cachedData)
            return
        }
    
        // Fetch paginated accounts from the database

        accounts, total, totalPages, err := paginatedAccounts(db, c)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching accounts"})
            return
        }
    
        // Prepare response data

        response := gin.H{
            "data":        accounts,
            "page":        c.Query("page"),
            "limit":       c.Query("limit"),
            "total":       total,
            "total_pages": totalPages,
        }
    
        // Cache the response data

        cacheInstance.Set(cacheKey, response, cache.DefaultExpiration)
    
        // Send response

        c.JSON(http.StatusOK, response)
    })
    

    r.Run(":8080")
}


    ////////////////////////////////////////// CACHING //////////////////////////////////////////
