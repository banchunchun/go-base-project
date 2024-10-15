package etcd

import (
	"com.banxiaoxiao.server/config"
	"com.banxiaoxiao.server/logger"
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Etcd struct {
	client   *clientv3.Client
	session  *concurrency.Session
	ctx      context.Context
	cancel   context.CancelFunc
	election *concurrency.Election
}

var etcdClient *Etcd

func NewEtcd() (*Etcd, error) {
	logger.Log().Infof("enter create etcd")
	if len(config.Cfg.Etcd.ServerList) <= 0 {
		logger.Log().Infof("leave create etcd with empty server list")
		return nil, nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	cli, err := clientv3.New(clientv3.Config{Endpoints: config.Cfg.Etcd.ServerList})
	if err != nil {
		logger.Log().Infof("failed to open etcd")
		return nil, err
	}
	session, err := concurrency.NewSession(cli)
	if err != nil {
		logger.Log().Infof("%s", err.Error())
		return nil, err
	}
	election := concurrency.NewElection(session, config.Cfg.Etcd.Leader)
	etcdClient = &Etcd{
		client:   cli,
		session:  session,
		ctx:      ctx,
		cancel:   cancel,
		election: election,
	}
	logger.Log().Infof("leave create etcd")
	return etcdClient, nil
}

func (c *Etcd) CloseEtcd() {
	c.cancel()
	if c.election != nil {
		c.election.Resign(c.ctx)
	}
	if c.session != nil {
		c.session.Close()
	}
	if c.client != nil {
		c.client.Close()
	}
}

func UpdateStorageToEtcd(v interface{}) {
	if etcdClient != nil {
		etcdClient.Put(v)
	}
}
