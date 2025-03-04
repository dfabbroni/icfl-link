package database

import (
	"fmt"

	"link/internal/config"
	"link/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Node{}, &models.Metadata{}, &models.Experiment{}, &models.ExperimentNode{})
	if err != nil {
		return nil, err
	}

	err = createDefaultUser(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDefaultUser(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("defaultpassword"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		defaultUser := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Approved: true,
		}

		result := db.Create(&defaultUser)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}