package main

import (
	"ecom/server/api"
	"ecom/server/handlers"
	"ecom/server/repos"
	"ecom/server/repos/postgres"
	"ecom/server/services"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	gormPg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(gormPg.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	db.AutoMigrate(
		&postgres.User{},
		&postgres.Address{},
		&postgres.Category{},
		&postgres.Product{},
		&postgres.ProductImage{},
		&postgres.Rating{},
		&postgres.Order{},
		&postgres.OrderItem{})
	repo := repos.IRepo(repos.NewPostgresDB(db))
	service := services.IService(services.NewService(repo))
	handlers := handlers.NewHandlers(service)
	app := api.NewApp(handlers)
	fmt.Println("🤠 server running at: ", os.Getenv("SRV_ADDR"))
	log.Fatal(app.Run(os.Getenv("SRV_ADDR")))

}
