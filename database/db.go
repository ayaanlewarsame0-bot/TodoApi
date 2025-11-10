package database

 import (
	
	"log"
	"os"
	"todo/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



var Db *gorm.DB

func InitDB() {
	// open the database
	connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
     fmt.Println("No DATABASE_URL found! Running locally?")
	 connStr = "user=postgres host=localhost password=Catarinax7 dbname=todo sslmode=disable"
	}
	var err error

	Db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// migrate the tables(object)
	err = Db.AutoMigrate(&models.User{}, &models.Todo{})
	if err != nil {
		log.Fatal("error migrating table", err)
	
}
}
