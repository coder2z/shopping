package discovery

import (
	"context"
	"fmt"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xjson"
	"github.com/coder2z/g-saber/xnet"
	"github.com/coder2z/g-saber/xstring"
	"go.etcd.io/etcd/clientv3"
	"shopping/constant"
	"shopping/utils"
	"sync"
	"time"
)

var (
	one sync.Once
	r   *RegDis
)

type RegDis struct {
	etcd    *clientv3.Client
	leaseId clientv3.LeaseID
	sync.WaitGroup
	closeCh    chan struct{}
	serverList []string
}

type ServerInfo struct {
	Ip   string
	Port int
}

func New(conf clientv3.Config, ch chan<- []ServerInfo) *RegDis {
	one.Do(func() {
		etcdC, err := clientv3.New(conf)
		if err != nil {
			panic(err)
		}
		r = &RegDis{etcd: etcdC, closeCh: make(chan struct{})}

		r.reg() //注册自己

		r.dis(ch) //发现别人
	})
	return r
}

func (r *RegDis) reg() {
	update := func() error {
		var (
			err  error
			step int
		)
		defer func() {
			if err != nil {
				utils.Log.Warn("etcd register error", err)
			}
		}()
		timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()
		ttl, err := r.etcd.Grant(timeout, 30)
		if err != nil {
			return err
		}

		step += 1
		ip, _, _ := xnet.GetLocalMainIP()
		info := ServerInfo{
			Ip:   ip,
			Port: xcfg.GetInt("server.port"),
		}
		_, err = r.etcd.Put(timeout, r.getKey(), xstring.Json(info), clientv3.WithLease(ttl.ID))
		if err == nil {
			r.leaseId = ttl.ID
		}
		return err
	}

	go func() {
		r.Add(1)
		defer r.Done()
		var err error
		err = update() // 先注册一次
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err == nil { // 注册成功则续租
					err = r.keepAliveOnce()
				}
				if err != nil { // 注册/续租失败则重新注册
					err = update()
				}
			case <-r.closeCh:
				r.unregister()
				return
			}
		}
	}()
}

func (r *RegDis) getKey() string {
	key := fmt.Sprintf("%s/%s", constant.EtcdKey, xstring.GenerateID())
	return key
}

func (r *RegDis) unregister() {
	key := r.getKey()
	if _, err := r.etcd.Delete(context.Background(), key); err != nil {
		utils.Log.Warn("unregister error", err)
	}
	_, _ = r.etcd.Revoke(context.Background(), r.leaseId)
}

func (r *RegDis) keepAliveOnce() error {
	_, err := r.etcd.KeepAliveOnce(context.Background(), r.leaseId)
	return err
}

func (r *RegDis) Close() {
	close(r.closeCh)
	r.Wait()
}

func (r *RegDis) dis(ch chan<- []ServerInfo) {
	update := func() []ServerInfo {
		resp, err := r.etcd.Get(context.Background(), constant.EtcdKey, clientv3.WithPrefix())
		if err != nil {
			utils.Log.Warn("etcd discovery watch", err)
			return nil
		}
		var i []ServerInfo
		for _, kv := range resp.Kvs {
			ins := ServerInfo{}
			if err = xjson.Unmarshal(kv.Value, &ins); err == nil {
				i = append(i, ins)
			}
		}
		return i
	}

	if i := update(); len(i) > 0 {
		ch <- i
	}

	eventCh := r.etcd.Watch(context.Background(), constant.EtcdKey, clientv3.WithPrefix())
	for range eventCh {
		ch <- update()
	}
	return
}
