package model

const (
	ActivityStatusNormal  = 0
	ActivityStatusDisable = 1
	ActivityStatusExpire  = 2
)

type Activity struct {
	ActivityId   int    `db:"id"`
	ActivityName string `db:"name"`
	ProductId    int    `db:"product_id"`
	StartTime    int64  `db:"start_time"`
	EndTime      int64  `db:"end_time"`
	Total        int    `db:"total"`
	Status       int    `db:"status"`

	StartTimeStr string
	EndTimeStr   string
	StatusStr    string
	Speed        int `db:"sec_speed"`
	BuyLimit     int `db:"buy_limit"`
}

type SecProductInfoConf struct {
	ProductId         int
	StartTime         int64
	EndTime           int64
	Status            int
	Total             int
	Left              int
	OnePersonBuyLimit int
	BuyRate           float64
	//每秒最多能卖多少个
	SoldMaxLimit int
}
