// Package store interacts with the SQL dialect and does data modeling
package store

import (
	"errors"
	"fmt"
	"github.com/ankurgel/duss/internal/duss/algo"
	"github.com/ankurgel/duss/internal/duss/models/auth"
	"github.com/ankurgel/duss/internal/duss/models/url"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// TODO: generalize to work for any dialect. Pick it from config
import _ "github.com/go-sql-driver/mysql" // This blank import is needed for gorm dialect to work

// Store represent a SQL binding for an adapter
type Store struct {
	Dialect string
	Db      *gorm.DB
}

// GormLogger is a custom logger for Gorm, making it use logrus.
type GormLogger struct{}

// Print handles log events from Gorm for the custom logger.
func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		log.WithFields(
			log.Fields{
				"module":  "gorm",
				"type":    "sql",
				"rows":    v[5],
				"src_ref": v[1],
				"values":  v[4],
			},
		).Info(v[3])
	case "log":
		log.WithFields(log.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}

// GetDSN returns Data Source Name for sql configuration
func (s *Store) GetDSN() string {
	config := viper.GetStringMapString("Mysql")
	host, username, password, database := config["host"], config["username"], config["password"], config["database"]
	if viper.GetString("Environment") == "development" {
		return fmt.Sprintf("%s:%s@/%s?parseTime=true", username, password, database)
	}
	return fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", username, password, host, database)
}

// EstablishConnection establishes connection of store with sql server
func (s *Store) EstablishConnection() {
	var err error
	s.Db, err = gorm.Open(s.Dialect, s.GetDSN())
	if err != nil {
		panic(fmt.Errorf("failed to connect to DB: %s", err))
	}
}

// SetupModels setups and migrates all the models
func (s *Store) SetupModels() {
	s.Db.AutoMigrate(&url.URL{})
	s.Db.AutoMigrate(&auth.User{})
}

// Close gracefully closes the store
func (s *Store) Close() {
	s.Db.Close()
}

// InitStore configures the store for connection, models,
// logging etc and returns instantiated store
func InitStore() *Store {
	s := &Store{Dialect: "mysql"}
	s.EstablishConnection()
	defer log.Info("Store configured successfully")
	s.Db.SetLogger(&GormLogger{})
	s.Db.LogMode(true)
	s.SetupModels()
	return s
}

// CreateByLongURL interacts with database to create short URL
// and returns URL object or error
func (s *Store) CreateByLongURL(longURL string, custom string) (*url.URL, error) {
	var u url.URL
	var shortURL *url.URL
	var err error
	var shortHash string
	if result := s.Db.Where("original = ?", longURL).First(&u); result.Error != nil{
		offset := 0
		if len(custom) > 0 {
			shortHash = custom
		} else {
			shortHash = algo.ComputeHash(longURL, offset)
		}

		shortURL, err = s.FindByShortURL(shortHash)
		// err will be nil if not found(happy), an object if found
		log.Error("-->", err)
		for err == nil && offset < viper.GetInt("MaxCollisionsAllowed") {
			log.Error("--> -->", offset, err)
			offset++
			shortHash = algo.ComputeHash(longURL, offset)
			shortURL, err = s.FindByShortURL(shortHash)
		}
		if shortURL == nil {
			short := url.URL{
				Short:      shortHash,
				Original:   longURL,
				Collisions: uint(offset),
			}
			if result := s.Db.Create(&short); result.Error != nil {
				return nil, errors.New("couldn't shorten. Something went wrong")
			}
			return &short, nil
		}
		return nil, errors.New("couldn't shorten. Out of lives")

	}
	return &u, nil
}

// FindByShortURL looks-up the store for given short url
// and returns URL object or error
func (s *Store) FindByShortURL(shortURL string) (*url.URL, error) {
	var u url.URL
	if result := s.Db.Where("short = ?", shortURL).First(&u); result.Error != nil{
		return nil, result.Error
	}
	return &u, nil
}

// CreateUser interacts with database to create a User
func (s *Store) CreateUser(user *auth.User) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("Password Encryption Failed")
	}
	user.Password = string(pass)
	token, err := user.CreateAPIToken()
	if err != nil {
		return "", errors.New("Token Creation Failed")
	}
	user.Token = token
	createdUser := s.Db.Create(user)
	if createdUser.Error != nil {
		return "", createdUser.Error
	}
	return token, err
}

// GetUserFromToken looksup the store to fetch a user associate with the token
func (s *Store) GetUserFromToken(token *auth.Token) (*auth.User, error) {
	user := &auth.User{}
	err := s.Db.Where("Email = ? AND Name = ?", token.Email, token.Name).First(user).Error
	return user, err
}