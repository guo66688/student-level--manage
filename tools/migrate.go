// tools\migrate.go
package main

import "student-level-manage/config"

// import (
// 	"fmt"

// 	"gorm.io/gorm"

// 	"student-level-manage/config"
// 	"student-level-manage/models"
// )

func main_6() {
	// 初始化数据库
	config.InitDB()
	// db := config.DB

	// migrateStudentClassNameToClassID(db)
}

// func migrateStudentClassNameToClassID(db *gorm.DB) {
// 	var students []models.Student
// 	if err := db.Find(&students).Error; err != nil {
// 		panic("❌ 查询学生失败: " + err.Error())
// 	}

// 	for _, stu := range students {
// 		if stu.ClassName == "" {
// 			continue // 跳过没有 ClassName 的
// 		}

// 		// 查找或创建班级
// 		var class models.Class
// 		if err := db.Where("name = ?", stu.ClassName).First(&class).Error; err != nil {
// 			if err == gorm.ErrRecordNotFound {
// 				class = models.Class{Name: stu.ClassName}
// 				if err := db.Create(&class).Error; err != nil {
// 					panic("❌ 创建班级失败: " + err.Error())
// 				}
// 				fmt.Printf("📌 创建新班级: %s (ID: %d)\n", class.Name, class.ID)
// 			} else {
// 				panic("❌ 查询班级失败: " + err.Error())
// 			}
// 		}

// 		// 更新学生记录
// 		if err := db.Model(&stu).Updates(map[string]interface{}{
// 			"ClassID":   class.ID,
// 			"ClassName": "", // 可选：清空旧字段
// 		}).Error; err != nil {
// 			panic(fmt.Sprintf("❌ 更新学生 [%s] 班级失败: %s", stu.Name, err.Error()))
// 		}
// 		fmt.Printf("✅ 学生 [%s] → 班级 [%s] (ClassID: %d)\n", stu.Name, class.Name, class.ID)
// 	}

// 	fmt.Println("🎉 班级字段迁移完成！")
// }
