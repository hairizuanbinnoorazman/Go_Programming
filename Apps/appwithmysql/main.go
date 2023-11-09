package main

import (
	"flag"
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
	var threadCount int
	var dataPoints int
	var hostname string
	flag.IntVar(&threadCount, "threads", 150, "Count of goroutine started to push data into database")
	flag.IntVar(&dataPoints, "points", 1000, "Number of data points to be pushed from a single goroutine")
	flag.StringVar(&hostname, "hostname", "localhost", "Hostname of database location")
	flag.Parse()

	fmt.Println("Start application")
	db, err := gorm.Open("mysql", fmt.Sprintf("username:password@tcp(%v:3306)/testmysql?charset=utf8&parseTime=True", hostname))
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

	for i := 1; i <= threadCount; i++ {
		wg.Add(1)
		go insertRecords(db, &wg, dataPoints)
	}

	for j := 1; j <= 10; j++ {
		go selectAll(db)
	}

	wg.Wait()

	defer db.Close()
}

func selectAll(db *gorm.DB) {

	for {
		var a []User
		result := db.Find(&a)
		if result.Error != nil {
			fmt.Printf("Error - unable to select db. %v", result.Error)
		}
		fmt.Printf("Number of records pulled: %v", len(a))
	}

}

func insertRecords(db *gorm.DB, wg *sync.WaitGroup, items int) {
	defer wg.Done()
	for i := 0; i < items; i++ {
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
