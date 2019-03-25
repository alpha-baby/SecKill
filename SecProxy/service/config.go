package service

import (
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

const (
	ProductStatusNormal       = 0
	ProductStatusSaleOut      = 1
	ProductStatusForceSaleOut = 2
)

type SecKillServer struct {
	RedisBlackConfig             RedisConfig
	RedisProxy2LayerConf         RedisConfig
	RedisLayer2ProxyConf         RedisConfig
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int

	EtcdConfig EtcdConfig

	SecProductInfoConfMap map[int]*SecProductInfoConf
	RWSecProductLock      sync.RWMutex

	LogPath  string
	LogLevel string

	CookieSecretKey string

	UserSecAccessLimit int

	ReferWhiteList       []string
	IpSecAccessLimit     int
	IpBlackMap           map[string]bool
	IdBlackMap           map[int]bool
	RWBlackLock          sync.RWMutex
	BlackRedisPool       *redis.Pool
	Proxy2LayerRedisPool *redis.Pool
	Layer2ProxyRedisPool *redis.Pool
	SecReqChan           chan *SecRequest
	SecReqChanSize       int

	secLimitMgr *SecLimitMgr
}

type RedisConfig struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}

type EtcdConfig struct {
	EtcdAddr         string
	EtcdTimeout      int
	EtcdSecKeyPrefix string
	EtcdProductKey   string
	EtcdBlackListKey string
}

type SecProductInfoConf struct {
	ProductId int
	StartTime int64
	EndTime   int64
	Status    int
	Total     int
	Left      int
}

type SecResult struct {
	ProductId int
	UserId    int
	Code      int
	Token     string
}

type SecRequest struct {
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
	CloseNotify   <-chan bool

	ResultChan chan *SecResult
}
