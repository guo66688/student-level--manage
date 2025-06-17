// contollers/class.go
package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"student-level-manage/config"
	"student-level-manage/models"

	"github.com/gin-gonic/gin"
)

// GetClasses godoc
// @Summary 获取班级列表
// @Tags classes
// @Accept json
// @Produce json
// @Param name query string false "班级名称（支持模糊匹配）" example("高一")
// @Success 200 {object} map[string]interface{} "返回班级列表和总数"
// @Router /classes [get]
func GetClasses(c *gin.Context) {
	var classes []models.Class
	var total int64

	name := c.DefaultQuery("name", "")

	// 如果无 name 查询条件，尝试从 Redis 缓存获取
	if name == "" {
		const cacheKey = "class:list"
		if val, err := config.Redis.Get(config.Ctx, cacheKey).Result(); err == nil {
			var cached []models.Class
			if json.Unmarshal([]byte(val), &cached) == nil {
				c.JSON(http.StatusOK, gin.H{
					"data":  cached,
					"total": len(cached),
				})
				return
			}
		}
	}

	// 构建查询条件
	db := config.DB.Model(&models.Class{})
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	if err := db.Find(&classes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}

	// 如果是无条件查询则缓存结果
	if name == "" {
		data, _ := json.Marshal(classes)
		config.Redis.Set(config.Ctx, "class:list", data, time.Hour)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  classes,
		"total": total,
	})
}

// 清除 Redis 缓存（用于新增/更新/删除时）
func invalidateClassCache() {
	config.Redis.Del(config.Ctx, "class:list")
}

// AddClass godoc
// @Summary 添加班级
// @Tags classes
// @Accept json
// @Produce json
// @Param data body models.Class true "班级信息" example({"name": "软件一班", "grade": 2023})
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /classes [post]
func AddClass(c *gin.Context) {
	var class models.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	if err := config.DB.Create(&class).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加失败"})
		return
	}
	invalidateClassCache()
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功", "class": class})
}

// UpdateClass godoc
// @Summary 更新班级
// @Tags classes
// @Accept json
// @Produce json
// @Param id path int true "班级ID" example(1)
// @Param data body models.Class true "班级信息" example({"name": "软件二班", "grade": 2023})
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /classes/{id} [put]
func UpdateClass(c *gin.Context) {
	id := c.Param("id")
	var class models.Class
	if err := config.DB.First(&class, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "班级不存在"})
		return
	}
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	if err := config.DB.Save(&class).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "更新失败"})
		return
	}
	invalidateClassCache()
	c.JSON(http.StatusOK, gin.H{"msg": "更新成功", "class": class})
}

// DeleteClass godoc
// @Summary 删除班级
// @Tags classes
// @Accept json
// @Produce json
// @Param id path int true "班级ID" example(1)
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /classes/{id} [delete]
func DeleteClass(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Class{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}
	invalidateClassCache()
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
