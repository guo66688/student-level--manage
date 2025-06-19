// controllers/chart.go
package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"student-level-manage/config"
	"student-level-manage/models"
	services "student-level-manage/services/charts"

	// Ensure the correct import
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ListCharts godoc
// @Summary 图表列表
// @Tags charts
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1) example(1)
// @Param size query int false "每页条数" default(10) example(10)
// @Param query query string false "查询关键词" example("成绩趋势")
// @Param type query string false "图表类型（可选）" example("score-trend")
// @Success 200 {object} map[string]interface{} "包含图表列表和总数"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /charts [get]
func ListCharts(c *gin.Context) {
	var charts []models.ChartData

	// 获取分页参数
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")
	query := c.DefaultQuery("query", "")
	typeParam := c.DefaultQuery("type", "")

	// 打印获取的参数
	fmt.Printf("Page: %s, Size: %s, Query: %s, Type: %s\n", page, size, query, typeParam)

	// 将字符串转换为整数
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "页码无效"})
		return
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "每页条数无效"})
		return
	}

	// 打印分页计算信息
	skip := int64((pageInt - 1) * sizeInt)
	limit := int64(sizeInt)
	fmt.Printf("Skip: %d, Limit: %d\n", skip, limit)

	// 构建查询条件
	filter := bson.M{}
	if query != "" {
		filter["meta.title"] = bson.M{"$regex": query, "$options": "i"}
	}
	if typeParam != "" {
		filter["type"] = typeParam
	}

	// 获取图表总数
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 打印正在查询图表总数
	fmt.Println("Fetching chart count...")
	count, err := config.MongoClient.Database("analysis").Collection("charts").CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败", "error": err.Error()})
		return
	}
	fmt.Printf("Chart count: %d\n", count)

	// 查询图表数据
	cursor, err := config.MongoClient.Database("analysis").Collection("charts").Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
		return
	}
	defer cursor.Close(ctx)

	// 解析查询结果
	fmt.Println("Parsing chart data...")
	if err := cursor.All(ctx, &charts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "解析失败"})
		return
	}

	// 打印解析后的图表数据
	fmt.Printf("Parsed %d charts\n", len(charts))

	// 构建返回的图表数据，符合新的格式：x 和 y 是独立的数组
	var result []map[string]interface{}
	for _, chart := range charts {
		// 构建图表的 values 数据
		var xValues []string
		var yValues []float64

		// 遍历 chart.Values 数组
		for _, value := range chart.Values {
			// 将 x 和 y 数据提取到独立的数组中
			xValues = append(xValues, value.X...)
			yValues = append(yValues, value.Y...)
		}

		// 将结果添加到返回数据中
		result = append(result, map[string]interface{}{
			"id":               chart.ID,
			"type":             chart.Type,
			"meta":             chart.Meta,
			"data_source_type": chart.DataSourceType,
			"values": map[string]interface{}{
				"x": xValues, // x 是字符串数组
				"y": yValues, // y 是数字数组
			},
		})
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": count,
	})
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

	// 直接使用 ShouldBindJSON 进行绑定
	if err := c.ShouldBindJSON(&chart); err != nil {
		fmt.Println("绑定错误:", err) // 打印具体的错误信息
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	// 打印接收到的图表数据
	fmt.Println("接收到的图表数据:", chart)

	// 遍历 chart.Values 数组，确保每个 ChartValues 中的 X 和 Y 数据有效
	for _, value := range chart.Values {
		// 检查 X 是否有效
		for _, x := range value.X {
			if x == "" {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "x 值不能为空"})
				return
			}
		}

		// 检查 Y 是否有效
		for _, y := range value.Y {
			if y <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "y 值必须是有效的数字"})
				return
			}
		}

		// 检查 X 和 Y 数组的长度是否匹配
		if len(value.X) != len(value.Y) {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "x 和 y 数组长度不匹配"})
			return
		}
	}

	// 存储数据到数据库
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := config.MongoClient.Database("analysis").Collection("charts").InsertOne(ctx, chart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "添加成功"})
}

// GetChartData godoc
// @Summary 获取图表数据
// @Tags charts
// @Accept json
// @Produce json
// @Param id query string true "图表ID" example("6852592bdc19f14fc8d8957c")
// @Param type query string true "图表类型" example("score-trend")
// @Param data_source_type query string true "数据源类型" example("static" 或 "dynamic")
// @Success 200 {object} map[string]interface{} "返回图表数据"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /charts/data [get]
func GetChartData(c *gin.Context) {
	fmt.Println("进入图表获取")

	// 获取查询参数
	chartID := c.DefaultQuery("id", "")
	chartType := c.DefaultQuery("type", "")
	dataSourceType := c.DefaultQuery("data_source_type", "dynamic")

	// 打印接收到的查询参数，检查是否正确传递
	fmt.Println("Received query params: id =", chartID, ", type =", chartType, ", data_source_type =", dataSourceType)

	// 如果没有传递图表ID，返回错误
	if chartID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "图表ID不能为空"})
		return
	}

	// 获取数据库连接（假设在 config 包中）
	db := config.MongoDB // 从 config 包中获取数据库连接

	// 使用工厂方法获取对应的生成器
	generator, err := services.ChartDataGeneratorFactory(chartType, dataSourceType, db)
	if err != nil {
		// 打印错误，帮助调试
		fmt.Println("Error in generator creation:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	// 打印确认生成器已成功创建
	fmt.Println("Generator created for chartType:", chartType, "and dataSourceType:", dataSourceType)

	// 根据 chartID 获取特定图表数据
	data, err := generator.GenerateDataByID(chartID)
	if err != nil {
		// 打印错误，帮助调试
		fmt.Println("Error fetching data for chartID:", chartID, "Error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成图表数据失败", "error": err.Error()})
		return
	}

	// 打印返回的数据，查看数据是否正确
	fmt.Println("Data fetched successfully for chartID:", chartID, "Data:", data)

	// 返回图表数据
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
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
