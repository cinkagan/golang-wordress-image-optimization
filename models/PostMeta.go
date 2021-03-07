package models

import (
	"keyholders/config"

	_ "github.com/go-sql-driver/mysql"
)

type PostMeta struct {
	MetaId    int32  `gorm:"column:meta_id"`
	PostId    string `gorm:"column:post_id"`
	MetaKey   string `gorm:"column:meta_key"`
	MetaValue string `gorm:"column:meta_value"`
}

func (PostMeta) TableName() string {
	return "wp_postmeta"
}

func GetAllListMeta() (_ []PostMeta, err error) {
	var posts []PostMeta
	if err = config.DB.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func GetSearchOnMeta(imgPath string) (_ []PostMeta, err error) {
	var posts []PostMeta
	if err = config.DB.Where("meta_value LIKE ?", "%"+imgPath+"%").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
