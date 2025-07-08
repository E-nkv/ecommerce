package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	pgGorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	"ecom/server/repos/postgres"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("err loading .env: ", err)
	}
	gofakeit.Seed(time.Now().UnixNano())
	gofakeit.LogLevel("error")
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(pgGorm.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	// Disable logging
	flushDB := false
	if len(os.Args) > 1 && os.Args[1] == "--flush" {
		flushDB = true
	}

	if flushDB {
		if err := flushDatabase(db); err != nil {
			log.Fatal("err flushing database: ", err)
		}
	}

	// Seed categories

	err = db.Transaction(func(tx *gorm.DB) error {
		categories := []postgres.Category{
			{Name: "Electronics"},
			{Name: "Food"},
			{Name: "Clothing"},
			{Name: "Books"},
			{Name: "Home"},
			{Name: "Toys"},
			{Name: "Sports"},
		}
		if err := tx.Create(&categories).Error; err != nil {
			return err
		}

		// Seed products
		var products []postgres.Product
		for range 5000 {
			cat := categories[rand.Intn(len(categories))]
			p := gofakeit.Product()
			products = append(products, postgres.Product{
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				CategoryID:  cat.ID,
			})
		}
		if err := tx.Create(&products).Error; err != nil {
			return err
		}

		// Seed users (including admins)
		var users []postgres.User
		for range 1000 {
			users = append(users, postgres.User{
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, 16),
				Name:     gofakeit.Name(),
				Role:     "user",
			})
		}
		// Add a few admin users
		users = append(users,
			postgres.User{Email: "enkv@gmail.com", Name: "nkv", Password: "nkv", Role: "admin"},
			postgres.User{Email: "perro@gmail.com", Name: "perro", Password: "perro", Role: "admin"},
		)
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error; err != nil {
			return err
		}

		// Fetch users from DB to get correct IDs
		var dbUsers []postgres.User
		if err := tx.Find(&dbUsers).Error; err != nil {
			return err
		}
		users = dbUsers

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
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&addresses).Error; err != nil {
			return err
		}

		// Fetch addresses from DB to get correct IDs
		var dbAddresses []postgres.Address
		if err := tx.Find(&dbAddresses).Error; err != nil {
			return err
		}
		addresses = dbAddresses

		// Seed orders and order items
		var orders []postgres.Order
		var orderItems []postgres.OrderItem
		for _, user := range users {
			numOrders := gofakeit.Number(1, 10)
			for range numOrders {
				// Get addresses for this user
				var userAddresses []postgres.Address
				for _, addr := range addresses {
					if addr.UserID == user.ID {
						userAddresses = append(userAddresses, addr)
					}
				}
				if len(userAddresses) == 0 {
					continue // skip if user has no addresses
				}
				addr := userAddresses[rand.Intn(len(userAddresses))]
				order := postgres.Order{
					UserID:        user.ID,
					AddressID:     addr.ID,
					Status:        gofakeit.RandomString([]string{"processing", "shipped", "delivered"}),
					TotalAmount:   0,
					PaymentMethod: gofakeit.RandomString([]string{"stripe", "paypal"}),
				}
				if err := tx.Create(&order).Error; err != nil {
					return err
				}
				numItems := gofakeit.Number(1, 5)
				var total float64
				for range numItems {
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
				// Update order total
				order.TotalAmount = total
				if err := tx.Model(&order).Update("total_amount", total).Error; err != nil {
					return err
				}
				orders = append(orders, order)
			}
		}
		if err := tx.CreateInBatches(&orderItems, 1000).Error; err != nil {
			return err
		}

		// Seed ratings
		var ratings []postgres.Rating
		for _, user := range users {
			for range gofakeit.Number(1, 5) {
				product := products[rand.Intn(len(products))]
				ratings = append(ratings, postgres.Rating{
					UserID:    user.ID,
					ProductID: product.ID,
					Score:     gofakeit.Number(1, 5),
					Comment:   gofakeit.Sentence(8),
				})
			}
		}
		if err := tx.Create(&ratings).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal("err seeding: ", err)
	}
	log.Println("Seeding complete!")
}
