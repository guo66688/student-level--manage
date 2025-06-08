// tools/update_class.go
package main

import (
	"fmt"
	"log"

	"student-level-manage/config"
	"student-level-manage/models"

	"gorm.io/gorm/clause"
)

func main_4() {
	// 初始化 DB
	config.InitDB()
	db := config.DB

	// 1) 确保已有的三个班级存在，不重复插入
	classes := []models.Class{
		{Name: "高一1班"},
		{Name: "高一2班"},
		{Name: "高一3班"},
	}
	if err := db.
		Clauses(clause.OnConflict{DoNothing: true}). // 如果已存在就跳过
		Create(&classes).
		Error; err != nil {
		log.Fatalf("插入 classes 失败：%v", err)
	}

	// 2) 拿到所有学生记录
	var students []models.Student
	if err := db.Find(&students).Error; err != nil {
		log.Fatalf("查询学生失败：%v", err)
	}

	// 3) 按 ID 排序后，分组更新 class_name
	for _, stu := range students {
		// 假设 ID 连续，从 1 开始，每 15 人为一班
		classNo := int((stu.ID-1)/15) + 1
		if classNo > 3 {
			classNo = 3 // 防止超出
		}
		newClass := fmt.Sprintf("高一%d班", classNo)

		if err := db.
			Model(&models.Student{}).
			Where("id = ?", stu.ID).
			Update("class_name", newClass).
			Error; err != nil {
			log.Printf("▶ 更新学生 ID=%d 班级失败：%v", stu.ID, err)
		}
	}

	log.Println("✅ 所有学生的班级字段已更新")
}
