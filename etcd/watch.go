package etcd

import (
	"com.banxiaoxiao.server/config"
	"com.banxiaoxiao.server/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (c *Etcd) Watch() {
	go func() {
		logger.Log().Infof("start watch on: %s", config.Cfg.Etcd.Alarm)
		watchChan := c.client.Watch(c.ctx, config.Cfg.Etcd.Alarm, clientv3.WithPrefix())
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				logger.Log().Infof("Event received! %s executed on %q with value %q\n", event.Type, event.Kv.Key, event.Kv.Value)
			}
		}
	}()
	go func() {
		logger.Log().Infof("start watch on: %s", config.Cfg.Etcd.Storage)
		watchChan := c.client.Watch(c.ctx, config.Cfg.Etcd.Storage, clientv3.WithPrefix())
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				logger.Log().Infof("Event received! %s executed on %q with value %q\n", event.Type, event.Kv.Key, event.Kv.Value)
			}
		}
	}()
}
