package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	FirstName  string `gorm:"type:varchar(100)"`
	LastName   string `gorm:"type:varchar(100)"`
	Occupation string `gorm:"type:varchar(50)"`
	Country    string `gorm:"type:varchar(30)"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func main() {
	fmt.Println("Start application")
	db, err := gorm.Open("mysql", "username:password@tcp(localhost:3306)/testmysql?charset=utf8&parseTime=True")
	db.DB().SetMaxOpenConns(100)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database %v", err))
	}
	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
	}
	db.AutoMigrate(&User{})
	fmt.Println("Managed to connect to the database")

	fmt.Println("Begin database entry loop")

	var wg sync.WaitGroup

	for i := 1; i <= 150; i++ {
		wg.Add(1)
		go insertRecords(db, &wg)
	}

	wg.Wait()

	defer db.Close()
}

func insertRecords(db *gorm.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		fmt.Printf("Iteration: %v\n", i)
		u := User{
			FirstName:  fmt.Sprintf("FName %v", i),
			LastName:   fmt.Sprintf("LName %v", i),
			Occupation: "testingoccupation",
			Country:    "singapore",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		db.Create(u)
	}
}
