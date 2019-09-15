package url

import (
	"github.com/jinzhu/gorm"
)

type Url struct {
	gorm.Model
	Short string `sql:"type:VARCHAR(7) CHARACTER SET utf8 COLLATE utf8_bin;unique"`
	//Short string `gorm:"unique;binary;not null"`
	Original string `gorm:"type:varchar(2048);index:orig"`
	Visits uint
}
