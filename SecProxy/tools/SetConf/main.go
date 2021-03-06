package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

const (
	EtcdKey = "/oldboy/backend/secskill/product"
)

type SecInfoConf struct {
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
}

func SetLogConfToEtcd() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")
	defer cli.Close()

	var SecInfoConfArr []SecInfoConf
	SecInfoConfArr = append(
		SecInfoConfArr,
		SecInfoConf{
			ID:135456454541621545,
			ProductId: 1029,
			StartTime: 1505008800,
			EndTime:   1505012400,
			Status:    0,
			Total:     1000,
			Left:      1000,
			OnePersonBuyLimit:3,
			BuyRate:0.5,
			//每秒最多能卖多少个
			SoldMaxLimit: 20,
		},
	)
	SecInfoConfArr = append(
		SecInfoConfArr,
		SecInfoConf{
			ID:135456454541621555,
			ProductId: 1027,
			StartTime: 1505008800,
			EndTime:   1505012400,
			Status:    0,
			Total:     2000,
			Left:      1000,
			OnePersonBuyLimit: 3,
			BuyRate:0.5,
			//每秒最多能卖多少个
			SoldMaxLimit: 20,
		},
	)

	data, err := json.Marshal(SecInfoConfArr)
	if err != nil {
		fmt.Println("json failed, ", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//cli.Delete(ctx, EtcdKey)
	//return
	_, err = cli.Put(ctx, EtcdKey, string(data))
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}

	//ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	//_, err = cli.Delete(ctx, EtcdKey)
	//cancelFunc()
	//if err != nil {
	//	log.Println("put faild :", err)
	//	return
	//}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}

func main() {
	SetLogConfToEtcd()
}
