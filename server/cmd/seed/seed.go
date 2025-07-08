package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	pgGorm "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ecom/server/repos/postgres"
)

func main() {
	gofakeit.Seed(time.Now().UnixNano())

	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(pgGorm.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Seed categories

	db.Transaction(func(tx *gorm.DB) error {
		categories := []postgres.Category{
			{Name: "Electronics"},
			{Name: "Food"},
			{Name: "Clothing"},
			{Name: "Books"},
			{Name: "Home & Kitchen"},
			{Name: "Toys"},
			{Name: "Sports"},
			{Name: "Beauty"},
			{Name: "Automotive"},
			{Name: "Garden"},
			{Name: "Health"},
			{Name: "Music"},
			{Name: "Office Supplies"},
			{Name: "Pet Supplies"},
			{Name: "Jewelry"},
			{Name: "Shoes"},
			{Name: "Baby"},
			{Name: "Movies"},
			{Name: "Video Games"},
			{Name: "Tools"},
		}
		if err := db.Create(&categories).Error; err != nil {
			return err
		}

		// Seed products
		var products []postgres.Product
		for i := 0; i < 5000; i++ {
			cat := categories[rand.Intn(len(categories))]
			p := gofakeit.Product()
			products = append(products, postgres.Product{
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				CategoryID:  cat.ID,
			})
		}
		if err := db.Create(&products).Error; err != nil {
			return err
		}

		// Seed users (including admins)
		var users []postgres.User
		for i := 0; i < 1000; i++ {
			users = append(users, postgres.User{
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, 16),
				Name:     gofakeit.Name(),
			})
		}
		// Add a few admin users
		for i := 0; i < 10; i++ {
			users = append(users, postgres.User{
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, 16),
				Name:     gofakeit.Name(),
				// Add an IsAdmin field if you have one
			})
		}
		if err := db.Create(&users).Error; err != nil {
			log.Fatalf("failed to seed users: %v", err)
		}

		// Seed addresses for users
		var addresses []postgres.Address
		for _, user := range users {
			for j := 0; j < gofakeit.Number(1, 3); j++ {
				addresses = append(addresses, postgres.Address{
					UserID:    user.ID,
					Line1:     gofakeit.Street(),
					City:      gofakeit.City(),
					State:     gofakeit.State(),
					ZipCode:   gofakeit.Zip(),
					Country:   gofakeit.Country(),
					IsDefault: j == 0,
				})
			}
		}
		if err := db.Create(&addresses).Error; err != nil {
			log.Fatalf("failed to seed addresses: %v", err)
		}

		// Seed orders and order items
		var orders []postgres.Order
		var orderItems []postgres.OrderItem
		for _, user := range users {
			numOrders := gofakeit.Number(1, 10)
			for j := 0; j < numOrders; j++ {
				addr := addresses[rand.Intn(len(addresses))]
				order := postgres.Order{
					UserID:        user.ID,
					AddressID:     addr.ID,
					Status:        gofakeit.RandomString([]string{"processing", "shipped", "delivered"}),
					TotalAmount:   0,
					PaymentMethod: gofakeit.RandomString([]string{"stripe", "paypal"}),
				}
				if err := db.Create(&order).Error; err != nil {
					log.Fatalf("failed to create order: %v", err)
				}
				numItems := gofakeit.Number(1, 5)
				var total float64
				for k := 0; k < numItems; k++ {
					product := products[rand.Intn(len(products))]
					quantity := gofakeit.Number(1, 3)
					item := postgres.OrderItem{
						OrderID:   order.ID,
						ProductID: product.ID,
						Quantity:  quantity,
						Price:     product.Price,
					}
					total += product.Price * float64(quantity)
					orderItems = append(orderItems, item)
				}
				order.TotalAmount = total
				if err := db.Save(&order).Error; err != nil {
					log.Fatalf("failed to update order total: %v", err)
				}
			}
		}
		if err := db.Create(&orderItems).Error; err != nil {
			log.Fatalf("failed to seed order items: %v", err)
		}

		// Seed ratings
		var ratings []postgres.Rating
		for _, user := range users {
			for i := 0; i < gofakeit.Number(1, 10); i++ {
				product := products[rand.Intn(len(products))]
				ratings = append(ratings, postgres.Rating{
					UserID:    user.ID,
					ProductID: product.ID,
					Score:     gofakeit.Number(1, 5),
					Comment:   gofakeit.Sentence(8),
				})
			}
		}
		if err := db.Create(&ratings).Error; err != nil {
			log.Fatalf("failed to seed ratings: %v", err)
		}

		return nil
	})
	log.Println("Seeding complete!")
}
