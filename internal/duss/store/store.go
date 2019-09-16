package store

import (
	"fmt"
	"github.com/ankurgel/duss/internal/duss/algo"
	"github.com/ankurgel/duss/internal/duss/models/url"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

import _ "github.com/go-sql-driver/mysql"

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

func (s *Store) GetDSN() string {
	config := viper.GetStringMapString("Mysql")
	host, port, username, password, database := config["host"], config["port"], config["username"], config["password"], config["database"]
	if viper.GetString("Environment") == "development" {
		return fmt.Sprintf("%s:%s@/%s?parseTime=true", username, password, database)
	} else {
		return fmt.Sprintf("%s:%s@%s:%s/%s?parseTime=true", username, password, host, port, database)
	}
}

func (s *Store) EstablishConnection() {
	var err error
	s.Db, err = gorm.Open(s.Dialect, s.GetDSN())
	if err != nil {
		panic(fmt.Errorf("Failed to connect to DB: %s\n", err))
	}
}

func (s *Store) SetupModels() {
	s.Db.AutoMigrate(&url.Url{})
}

func (s *Store) Close() {
	s.Db.Close()
}

func InitStore() *Store {
	s := &Store{Dialect: "mysql"}
	s.EstablishConnection()
	defer log.Info("Store configured successfully")
	s.Db.SetLogger(&GormLogger{})
	s.Db.LogMode(true)
	s.SetupModels()
	return s
}

func (s *Store) CreateByLongUrl(longUrl string) (*url.Url, error) {
	var u url.Url
	var shortUrl *url.Url
	var err error
	if result := s.Db.Where("original = ?", longUrl).First(&u); result.Error != nil{
		offset := 0
		shortHash := algo.ComputeHash(longUrl, offset)

		shortUrl, err = s.FindByShortUrl(shortHash)
		// err will be nil if not found(happy), an object if found
		log.Error("-->", err)
		for err == nil && offset < 5 {
			log.Error("--> -->", offset, err)
			offset++
			shortHash = algo.ComputeHash(longUrl, offset)
			shortUrl, err = s.FindByShortUrl(shortHash)
		}
		if shortUrl == nil {
			result := url.Url{
				Short: shortHash,
				Original: longUrl,
				Collisions: uint(offset),
			}
			s.Db.Create(&result)
			return &result, nil
		} else {
			return nil, errors.New("Couldn't shorten. Out of lives")
		}

	} else {
		return &u, nil
	}



	return nil, errors.New("Cannot shorten. Out of lives.")
}

func (s *Store) FindByShortUrl(shortUrl string) (*url.Url, error) {
	var u url.Url
	if result := s.Db.Where("short = ?", shortUrl).First(&u); result.Error != nil{
		return nil, result.Error
	}
	return &u, nil
}