// tools\seed.go
package main

import (
	"fmt"
	"time"

	"student-level-manage/config"
	"student-level-manage/models"

	"gorm.io/gorm/clause"
)

func main() {
	// 1. 初始化数据库
	config.InitDB()
	db := config.DB

	// 2. 确保三张班级已经存在（无冲突则插入）
	classes := []models.Class{
		{Name: "高一1班"},
		{Name: "高一2班"},
		{Name: "高一3班"},
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&classes).Error; err != nil {
		panic("❌ 插入班级失败: " + err.Error())
	}

	// 3. 读出这三条班级到 map[className]Class
	var allClasses []models.Class
	if err := db.Where("name IN ?", []string{"高一1班", "高一2班", "高一3班"}).Find(&allClasses).Error; err != nil {
		panic("❌ 读取班级失败: " + err.Error())
	}
	classMap := make(map[string]models.Class, len(allClasses))
	for _, cls := range allClasses {
		classMap[cls.Name] = cls
	}

	// 4. 读已有学生，避免重复
	var existing []models.Student
	db.Find(&existing)
	existMap := make(map[string]bool, len(existing))
	for _, s := range existing {
		existMap[s.StudentNo] = true
	}

	// 5. 插入学生
	names := []string{
		"罗惠萍", "罗汉平", "肖利萍", "胡建国", "苏汉明", "萧汉文", "蒋小平", "蒲志辉", "蔡泽强", "袁永强",
		"蒋玉珍", "金琴娥", "荆红梅", "景炳芳", "鞠华美", "康纪君", "冷慧", "李保翠", "李昌惠", "李付玉",
		"李洪芝", "李腊琴", "李明法", "李习连", "李玉花",
		"黄伟林", "黄桂潮", "黎淇棠", "龙小文", "陈彩香", "陈华均", "周建军", "武秀俊", "包广杰", "张淑梅",
		"孙淑霞", "殷丽坤", "周家秀", "崔树森", "扈玉柱",
	}

	var toCreate []models.Student
	for i, name := range names {
		studentNo := fmt.Sprintf("202300%02d", i+1)
		if existMap[studentNo] {
			continue
		}
		className := fmt.Sprintf("高一%d班", i/15+1)
		cls, ok := classMap[className]
		if !ok {
			panic("找不到班级：" + className)
		}

		gender := "男"
		if i%2 == 0 {
			gender = "女"
		}
		birth := time.Date(2005, time.January, (i%28)+1, 0, 0, 0, 0, time.Local)

		toCreate = append(toCreate, models.Student{
			StudentNo: studentNo,
			Name:      name,
			Gender:    gender,
			BirthDate: birth.Format("2006-01-02"),
			ClassID:   cls.ID,
		})
	}

	// 插入学生数据，不退出
	if len(toCreate) > 0 {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&toCreate).Error; err != nil {
			panic("❌ 插入学生失败: " + err.Error())
		}
		fmt.Printf("✅ 新增 %d 名学生\n", len(toCreate))
	} else {
		fmt.Println("✅ 所有学生已存在，跳过插入")
	}

	// 6. 插入三门课程（同理去重）
	courses := []models.Course{
		{CourseCode: "MATH101", CourseName: "数学", Credit: 3, TeacherID: 1},
		{CourseCode: "CHIN101", CourseName: "语文", Credit: 3, TeacherID: 2},
		{CourseCode: "ENGL101", CourseName: "英语", Credit: 3, TeacherID: 3},
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&courses).Error; err != nil {
		panic("❌ 插入课程失败: " + err.Error())
	}
	fmt.Println("✅ 课程表已更新")

	// 7. 插入成绩（如果有重复，则跳过）
	var allStudents []models.Student
	var allCourses []models.Course
	db.Find(&allStudents)
	db.Find(&allCourses)

	examDates := []string{
		"2025-05-15", "2025-12-15", "2025-01-15",
		"2025-02-15", "2025-03-15", "2025-04-15",
		"2025-06-15", "2025-09-15", "2025-10-15", "2025-11-15",
	}

	var scores []models.Score
	for _, stu := range allStudents {
		for _, course := range allCourses {
			for _, date := range examDates {
				// 检查该学生-课程-日期是否已经有成绩记录
				var existingScore models.Score
				db.Where("student_id = ? AND course_id = ? AND exam_date = ?", stu.ID, course.ID, date).First(&existingScore)
				if existingScore.ID != 0 {
					continue // 如果存在成绩，跳过
				}

				// 生成成绩
				val := 60 + float64((stu.ID+course.ID)%41)
				scores = append(scores, models.Score{
					StudentID: stu.ID,
					CourseID:  course.ID,
					ExamDate:  date,
					Score:     val,
				})
			}
		}
	}

	// 插入成绩
	if len(scores) > 0 {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&scores).Error; err != nil {
			panic("❌ 插入成绩失败: " + err.Error())
		}
		fmt.Printf("✅ 尝试插入 %d 条新的成绩记录（冲突自动跳过）\n", len(scores))
	} else {
		fmt.Println("✅ 没有新成绩需要插入")
	}
}
