package service

import (
	"github.com/garyburd/redigo/redis"
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

var (
	seclayerContext = &SecLayerContext{}
)

type SecLayerConf struct {
	Proxy2LayerRedis             RedisConfig
	Layer2ProxyRedis             RedisConfig
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int
	HandleUserGoroutineNum       int
	Read2HandleChanSize          int
	Handle2WriteChanSize         int
	MaxRequestWaitTimeout        int64
	SendToWriteChanTimeout       int
	SendToHandleChanTimeout      int

	EtcdConfig EtcdConfig

	LogLevel string
	LogPath  string

	SecProductInfoMap map[int64]*SecProductInfoConf

	TokenPasswd string
}

type RedisConfig struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
	RedisQueueName   string
}

type EtcdConfig struct {
	EtcdAddr          string
	Timeout           int
	EtcdSecKeyPrefix  string
	EtcdSecProductKey string
}

type SecLayerContext struct {
	Proxy2LayerRedisPool *redis.Pool
	Layer2ProxyRedisPool *redis.Pool
	EtcdClient           *clientv3.Client
	RWSecProductLock     sync.RWMutex

	SecLayerConf *SecLayerConf

	waitGroup        sync.WaitGroup
	Read2HandleChan  chan *SecRequest
	Handle2WriteChan chan *SecResponse
	HistoryMap       map[int]*UserBuyHistory
	HistoryMapLock   sync.Mutex

	//商品的计数
	ProductCountMgr *ProductCountMgr
}

type SecProductInfoConf struct {
	ID int64
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
	// 限速控制
	secLimit *SecLimit

}

type SecRequest struct {
	ID int64
	ProductId     int
	Source        string
	AuthCode      string
	SecTime       string
	Nance         string
	UserId        int
	UserAuthSign  string
	AccessTime    time.Time
	ClientAddr    string
	ClientRefence string
	//CloseNotify   <-chan bool

	//ResultChan chan *SecResult
}

type SecResponse struct {
	ID int64
	ProductId int
	UserId    int
	Token     string
	Code      int
	TokenTime int64
}

type SecResult struct {
	ID int64
	ProductId int
	UserId    int
	Code      int
	Token     string
}
