package redisop

import (
	"context"
	"strconv"
	"student-level-manage/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

const RankKey = "score_rank"

// 写入排行榜
func UpdateStudentRank(studentID uint, avgScore float64) error {
	return config.Redis.ZAdd(ctx, RankKey, redis.Z{
		Score:  avgScore,
		Member: studentID,
	}).Err()
}

// 获取前 N 名
func GetTopN(n int64) ([]redis.Z, error) {
	return config.Redis.ZRevRangeWithScores(ctx, RankKey, 0, n-1).Result()
}

// 可选：移除
func RemoveStudent(studentID uint) error {
	return config.Redis.ZRem(ctx, RankKey, studentID).Err()
}

func RefreshAllRanksFromDB() error {
	var results []struct {
		StudentID uint
		AvgScore  float64
	}

	// 从 MySQL 聚合平均成绩
	err := config.DB.
		Table("scores").
		Select("student_id, AVG(score) as avg_score").
		Group("student_id").
		Scan(&results).Error
	if err != nil {
		return err
	}

	pipe := config.Redis.TxPipeline()
	pipe.Del(ctx, RankKey) // 清空旧排行

	for _, item := range results {
		pipe.ZAdd(ctx, RankKey, redis.Z{
			Score:  item.AvgScore,
			Member: item.StudentID,
		})
	}
	pipe.Expire(ctx, RankKey, time.Hour)
	_, err = pipe.Exec(ctx)
	return err
}

func GetStudentRank(studentID uint) (int64, error) {
	return config.Redis.ZRevRank(ctx, RankKey, strconv.Itoa(int(studentID))).Result()
}
