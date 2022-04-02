package util

import (
	"github.com/mediocregopher/radix/v3"
)

type Pool interface {
	Do(a radix.Action) error
	Close() error
}

func ConnectToRedis(config DatabaseConfig) (Pool, error) {
	if config.ClusterModeEnabled {
		poolFunc := func(network, addr string) (radix.Client, error) {
			return radix.NewPool(network, addr, config.PoolSize)
		}
		return radix.NewCluster(config.Addresses, radix.ClusterPoolFunc(poolFunc))
	} else {
		return radix.NewPool("tcp", config.Addresses[0], config.PoolSize)
	}
}

func CloseRedis(config DatabaseConfig, pool Pool) {
	pool.Close()
}
