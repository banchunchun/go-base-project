package etcd

import (
	"com.banxiaoxiao.server/config"
	"com.banxiaoxiao.server/logger"
)

func (c *Etcd) Leader() error {
	logger.Log().Infof("enter leader election")
	if config.Cfg.Etcd.Election {
		if err := c.election.Campaign(c.ctx, "e"); err != nil {
			logger.Log().Errorf("%s", err.Error())
			return err
		}
	}
	logger.Log().Infof("leave leader election")
	return nil
}
