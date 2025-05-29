// services/chart.go
package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"student-level-manage/config"
	"student-level-manage/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetChartData(chartType string) ([]models.ChartData, error) {
	collection := config.MongoDB.Collection("charts")
	cur, err := collection.Find(context.Background(), bson.M{"type": chartType})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var results []models.ChartData
	for cur.Next(context.Background()) {
		var item models.ChartData
		if err := cur.Decode(&item); err == nil {
			results = append(results, item)
		}
	}
	return results, nil
}

// 保存或更新图表数据（按 type 去重）
func SaveOrUpdateChartData(ctx context.Context, data models.ChartData) {
	filter := bson.M{"type": data.Type}
	update := bson.M{
		"$set": bson.M{
			"type":   data.Type,
			"meta":   data.Meta,
			"values": data.Values,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := config.MongoClient.
		Database("analysis").
		Collection("charts").
		UpdateOne(ctx, filter, update, opts)

	if err != nil {
		log.Println("❌ MongoDB 更新图表数据失败:", err)
	}
}

func UpdatePassRateChart(ctx context.Context) {
	// 不再需要自己创建 ctx
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// 查询每门课的通过率（分数 >= 60）
	type Item struct {
		CourseID uint   `json:"course_id" bson:"course_id"`
		Passed   int    `json:"passed" bson:"passed"`
		Total    int    `json:"total" bson:"total"`
		Rate     string `json:"rate" bson:"rate"`
	}
	var items []Item

	// 查询所有课程 ID
	var courseIDs []uint
	config.DB.Table("scores").Select("DISTINCT course_id").Scan(&courseIDs)

	for _, cid := range courseIDs {
		var total, passed int64
		config.DB.Model(&models.Score{}).Where("course_id = ?", cid).Count(&total)
		config.DB.Model(&models.Score{}).Where("course_id = ? AND score >= 60", cid).Count(&passed)

		rate := "0%"
		if total > 0 {
			rate = fmt.Sprintf("%.2f%%", float64(passed)/float64(total)*100)
		}
		items = append(items, Item{
			CourseID: cid,
			Passed:   int(passed),
			Total:    int(total),
			Rate:     rate,
		})
	}

	SaveOrUpdateChartData(ctx, models.ChartData{
		Type:   "pass_rate",
		Meta:   map[string]interface{}{},
		Values: items,
	})
}

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
