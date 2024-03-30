package intergration

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"project0/interactive/intergration/startup"
	"project0/interactive/repository/dao"
	"testing"
	"time"
)

/*
User: society-programmer
Date: 2024/3/15  周五
Time: 16:17
*/



//go:embed init.sql
var initSQL string

func TestGenSQL(t *testing.T)  {
	file,err := os.OpenFile("data.sql",
		os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0666)
	assert.NoError(t, err)

	defer file.Close()

	_, err  = file.WriteString(initSQL)

	require.NoError(t, err)



}

func TestGenData(t *testing.T) {
	// 这个是批量插入，数据量不是特别大的时候，可以用这个
	// GenData 要比 GenSQL 慢
	// 你根据自己的需要调整批次，和每个批次大小
	db := startup.InitDB()
	// 这个为 true，只会输出 SQL，但是不会执行，也不会报错
	// db.DryRun = true
	// 1000 批
	for i := 0; i < 10; i++ {
		// 每次 100 条
		// 你可以考虑直接用 CreateInBatches，GORM 帮你分批次
		// 我自己分是为了控制内存消耗
		const batchSize = 100
		data := make([]dao.Interactive, 0, batchSize)
		now := time.Now().UnixMilli()
		for j := 0; j < batchSize; j++ {
			data = append(data, dao.Interactive{
				Biz:        "test",
				BizId:      int64(i*batchSize + j + 1),
				ReadCnt:    rand.Int63(),
				LikeCnt:    rand.Int63(),
				CollectCnt: rand.Int63(),
				Utime:      now,
				Ctime:      now,
			})
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			err := tx.Create(data).Error
			require.NoError(t, err)
			return err
		})
		require.NoError(t, err)
	}
}