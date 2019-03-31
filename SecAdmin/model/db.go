package model

import (
	"github.com/jmoiron/sqlx"
	"go.etcd.io/etcd/clientv3"
)

var (
	Db *sqlx.DB
	EtcdClient *clientv3.Client
	EtcdPrefix string
	EtcdProductKey string
)

func InitDB(db *sqlx.DB)  {
	Db = db
	return
}

func InitEtcd(c *clientv3.Client, prefix string, productKey string) {
	EtcdClient = c
	EtcdPrefix = prefix
	EtcdProductKey = productKey
	return
}