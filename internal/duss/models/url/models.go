// Package url contains the URL data model
package url

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// URL data model
type URL struct {
	gorm.Model
	Short string `sql:"type:VARCHAR(7) CHARACTER SET utf8 COLLATE utf8_bin;unique"`
	//Short string `gorm:"unique;binary;not null"`
	Original string `gorm:"type:varchar(2048);index:orig"`
	Visits uint
	Collisions uint
}

// ShortURL gives http URL string form for a short slug
func (u *URL) ShortURL() string {
	return viper.GetString("BaseURL") +  u.Short
}