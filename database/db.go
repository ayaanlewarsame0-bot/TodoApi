package database

 import (
	"todo/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



var Db *gorm.DB

func InitDB() {
	// open the database
	connStr := "user=postgres host=localhost password=Catarinax7 dbname=todo sslmode=disable"

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
