package etcd

import (
	"com.banxiaoxiao.server/config"
	"com.banxiaoxiao.server/logger"
	"com.banxiaoxiao.server/util"
)

func (c *Etcd) Put(v interface{}) error {
	r, err := c.client.Put(c.ctx, config.Cfg.Etcd.Storage, util.ToString(v))
	logger.Log().Infof("put to storage %+v", r)
	return err
}
