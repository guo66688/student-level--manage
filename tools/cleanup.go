package main

import (
	"fmt"

	"student-level-manage/config"
	"student-level-manage/models"
)

func main_5() {
	config.InitDB()
	db := config.DB

	// ✅ 删除所有表（自动检测表依赖顺序）
	err := db.Migrator().DropTable(
		&models.Score{},
		&models.Student{},
		&models.Course{},
		&models.Class{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Permission{},
		&models.Role{},
		&models.User{},
	)
	if err != nil {
		panic("❌ 删除表失败: " + err.Error())
	}
	fmt.Println("✅ 所有表已清除")

	// ✅ 重新建表（基于更新后的模型）
	err = db.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Course{},
		&models.Score{},
		&models.Permission{},
		&models.Role{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.Class{},
	)
	if err != nil {
		panic("❌ 重新建表失败: " + err.Error())
	}
	fmt.Println("✅ 表结构已重建")
}
