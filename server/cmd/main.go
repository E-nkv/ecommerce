package main

import (
	"context"
	"ecom/server/api"
	"ecom/server/handlers"
	"ecom/server/repos"
	"ecom/server/repos/products"
	productsService "ecom/server/services/products"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dburl := os.Getenv("DB_URL")
	db, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	var productRepo repos.IProductRepo = products.NewProductRepo(db)
	var productService *productsService.ProductService = productsService.NewService(productRepo)

	handlers := handlers.NewHandlers(productService)
	app := api.NewApp(handlers)
	fmt.Println("🤠 server running at: ", os.Getenv("SRV_ADDR"))
	log.Fatal(app.Run(os.Getenv("SRV_ADDR")))

}
