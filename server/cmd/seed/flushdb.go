package main

import (
	"gorm.io/gorm"
)

func flushDatabase(db *gorm.DB) error {
	err := db.Exec(`
        TRUNCATE TABLE 
            order_items, 
            orders, 
            ratings, 
            addresses, 
            users, 
            products, 
            categories 
        RESTART IDENTITY CASCADE
    	`).Error
	return err
}
