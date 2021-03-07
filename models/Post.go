package models

import (
	"keyholders/config"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	ID        int32  `gorm:"column:ID"`
	PostTitle string `gorm:"column:post_title"`
	PostName  string `gorm:"column:post_name"`
	GUID      string `gorm:"column:guid"`
}

func (Post) TableName() string {
	return "wp_posts"
}

func GetAllList() (_ []Post, err error) {
	var posts []Post
	if err = config.DB.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func GetSearchOnPost(imgPath string) (_ []Post, err error) {
	var posts []Post
	if err = config.DB.Where("guid LIKE ?", "%"+imgPath+"%").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
