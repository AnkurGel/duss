package store

import (
	"fmt"
	"github.com/ankurgel/duss/internal/duss/models/url"
	"github.com/jinzhu/gorm"
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
		return fmt.Sprintf("%s:%s@/%s", username, password, database)
	} else {
		return fmt.Sprintf("%s:%s@%s:%s/%s", username, password, host, port, database)
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
	return s
}

func (s *Store) FindBySlug(slug string) (*url.Url, error) {
	url := url.Url{}
	if result := s.Db.Where("short = ?", slug).First(&url); result.Error != nil{
		return nil, result.Error
	}
	return &url, nil
}