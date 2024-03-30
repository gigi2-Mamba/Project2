package dao

import "gorm.io/gorm"

/*
User: society-programmer
Date: 2024/3/6  周三
Time: 10:34
*/
// 不需要。 手动维护表不需要
func InitTables(db *gorm.DB) error  {


	return db.AutoMigrate(
		&Interactive{},
		&UserLikeBiz{},
		&UserCollectBiz{},
	)
}