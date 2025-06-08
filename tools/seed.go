// tools/seed.go
// package main

// import (
// 	"fmt"
// 	"time"

// 	"student-level-manage/config"
// 	"student-level-manage/models"
// )

// func main() {
// 	// 1. 先初始化数据库
// 	config.InitDB()
// 	db := config.DB

// 	// 2. 学生名单，45 人，分 3 班
// 	names := []string{
// 		"罗惠萍", "罗汉平", "肖利萍", "胡建国", "苏汉明", "萧汉文", "蒋小平", "蒲志辉", "蔡泽强", "袁永强",
// 		"蒋玉珍", "金琴娥", "荆红梅", "景炳芳", "鞠华美", "康纪君", "冷慧", "李保翠", "李昌惠", "李付玉",
// 		"李洪芝", "李腊琴", "李明法", "李习连", "李玉花",
// 		"黄伟林", "黄桂潮", "黎淇棠", "龙小文", "陈彩香", "陈华均", "周建军", "武秀俊", "包广杰", "张淑梅",
// 		"孙淑霞", "殷丽坤", "周家秀", "崔树森", "扈玉柱",
// 	}

// 	var students []models.Student
// 	for i, name := range names {
// 		classNo := i/15 + 1
// 		className := fmt.Sprintf("高一%d班", classNo)
// 		gender := "男"
// 		if i%2 == 0 {
// 			gender = "女"
// 		}
// 		studentNo := fmt.Sprintf("202300%02d", i+1)
// 		birth := time.Date(2005, time.January, (i%28)+1, 0, 0, 0, 0, time.Local)

// 		students = append(students, models.Student{
// 			Name:      name,
// 			StudentNo: studentNo,
// 			Gender:    gender,
// 			BirthDate: birth.Format("2006-01-02"), // 如果模型字段是 string
// 			ClassName: className,
// 		})
// 	}

// 	if err := db.Create(&students).Error; err != nil {
// 		panic("❌ 插入学生失败: " + err.Error())
// 	}
// 	fmt.Printf("✅ 插入 %d 名学生成功\n", len(students))

// 	// 3. 三门课程
// 	courses := []models.Course{
// 		{CourseCode: "MATH101", CourseName: "数学", Credit: 3, TeacherID: 1},
// 		{CourseCode: "CHIN101", CourseName: "语文", Credit: 3, TeacherID: 2},
// 		{CourseCode: "ENGL101", CourseName: "英语", Credit: 3, TeacherID: 3},
// 	}
// 	if err := db.Create(&courses).Error; err != nil {
// 		panic("❌ 插入课程失败: " + err.Error())
// 	}
// 	fmt.Printf("✅ 插入 %d 门课程成功\n", len(courses))

//		// 4. 期中和期末成绩
//		var scores []models.Score
//		examDates := []string{"2025-05-15", "2025-12-15"}
//		for _, stu := range students {
//			for _, course := range courses {
//				for _, date := range examDates {
//					// 简单用学号+课程序号生成一个 60-100 范围内的分数
//					val := 60 + float64((stu.ID+course.ID)%41)
//					scores = append(scores, models.Score{
//						StudentID: stu.ID,
//						CourseID:  course.ID,
//						ExamDate:  date,
//						Score:     val,
//					})
//				}
//			}
//		}
//		if err := db.Create(&scores).Error; err != nil {
//			panic("❌ 插入成绩失败: " + err.Error())
//		}
//		fmt.Printf("✅ 插入 %d 条成绩记录成功\n", len(scores))
//	}
package main

import (
	"fmt"
	"log"

	"student-level-manage/config"
	"student-level-manage/models"

	"gorm.io/gorm/clause"
)

func main() {
	config.InitDB()
	db := config.DB

	// 1) 确保三张班级已经存在
	classes := []models.Class{
		{Name: "高一1班"},
		{Name: "高一2班"},
		{Name: "高一3班"},
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&classes).Error; err != nil {
		log.Fatalf("插入 classes 失败：%v", err)
	}

	// 2) 读取所有学生
	var students []models.Student
	if err := db.Find(&students).Error; err != nil {
		log.Fatalf("查学生失败：%v", err)
	}

	// 3) 遍历并更新 class_name
	for _, stu := range students {
		// 按照你的逻辑：每 15 人一班
		// 假设 stu.ID 从 1 开始且是连续的
		classNo := int((stu.ID-1)/15) + 1
		className := fmt.Sprintf("高一%d班", classNo)

		if err := db.Model(&models.Student{}).
			Where("id = ?", stu.ID).
			Update("class_name", className).
			Error; err != nil {
			log.Printf("更新学生 %d 班级失败：%v", stu.ID, err)
		}
	}
	log.Println("✅ 所有学生的班级字段已更新")
}
