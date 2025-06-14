// controllers/chart.go
package controllers

import (
	"context"
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ListCharts godoc
// @Summary 图表列表
// @Tags charts
// @Accept json
// @Produce json
// @Param type query string false "图表类型（可选）" example("score-trend")
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /charts [get]
// 图表列表接口（支持按类型筛选）
func ListCharts(c *gin.Context) {
	typeParam := c.Query("type")
	filter := bson.M{}
	if typeParam != "" {
		filter["type"] = typeParam
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := config.MongoClient.Database("analysis").Collection("charts").Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	defer cursor.Close(ctx)

	var charts []models.ChartData
	if err := cursor.All(ctx, &charts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "解析失败"})
		return
	}
	c.JSON(http.StatusOK, charts)
}

// AddChart godoc
// @Summary 新增图表
// @Tags charts
// @Accept json
// @Produce json
// @Param data body models.ChartData true "图表数据" example({"type":"score-trend","meta":{"title":"成绩趋势"},"values":[{"x":"2024-01","y":85.3}]})
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /charts [post]
// 图表新增接口
func AddChart(c *gin.Context) {
	var chart models.ChartData
	if err := c.ShouldBindJSON(&chart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := config.MongoClient.Database("analysis").Collection("charts").InsertOne(ctx, chart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功"})
}

// UpdateChart godoc
// @Summary 更新图表
// @Tags charts
// @Accept json
// @Produce json
// @Param id path string true "图表ID" example("60f7d2f5e13f1e3c48a2a1a2")
// @Param data body models.ChartData true "图表数据" example({"type":"pass-rate","meta":{"title":"通过率"},"values":[{"x":"课程1","y":92.5}]})
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /charts/{id} [put]
// 图表编辑接口（根据 ID 更新）
func UpdateChart(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的 ID"})
		return
	}

	var chart models.ChartData
	if err := c.ShouldBindJSON(&chart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = config.MongoClient.Database("analysis").Collection("charts").UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"type":   chart.Type,
			"meta":   chart.Meta,
			"values": chart.Values,
		}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

// DeleteChart godoc
// @Summary 删除图表
// @Tags charts
// @Accept json
// @Produce json
// @Param id path string true "图表ID" example("60f7d2f5e13f1e3c48a2a1a2")
// @Success 200 {object} map[string]interface{} "返回信息"
// @Router /charts/{id} [delete]
// 删除图表接口
func DeleteChart(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := config.MongoClient.
		Database("analysis").
		Collection("charts").
		DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}
	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "未找到图表"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
