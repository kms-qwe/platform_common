package redis

import (
	"context"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/kms-qwe/platform_common/pkg/client/cache"
)

type handler func(ctx context.Context, conn redis.Conn) error

type client struct {
	pool              *redis.Pool
	connectionTimeout time.Duration
}

// NewClient create new RedishCache client
func NewClient(pool *redis.Pool, connectionTimeout time.Duration) cache.RedisCache {
	return &client{
		pool:              pool,
		connectionTimeout: connectionTimeout,
	}
}

// getConnect returns connect to redis
func (c *client) getConnect(ctx context.Context) (redis.Conn, error) {
	ctx, cancel := context.WithTimeout(ctx, c.connectionTimeout)
	defer cancel()

	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		log.Printf("failed to get redis connection: %v\n", err)

		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}

// execute gets conneetct to redis and executes hanlder func
func (c *client) execute(ctx context.Context, handler handler) error {
	conn, err := c.getConnect(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("failed to close redis connection: %v\n", err)
		}
	}()

	err = handler(ctx, conn)
	if err != nil {
		return err
	}

	return nil
}

// Ping pingsg redis
func (c *client) Ping(ctx context.Context) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		_, err := conn.Do("PING")
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Set sets ket to value in redis
func (c *client) Set(ctx context.Context, key string, value interface{}) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		_, err := conn.Do("SET", redis.Args{key}.Add(value)...)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Get gets values by key in redis
func (c *client) Get(ctx context.Context, key string) (interface{}, error) {
	var value interface{}
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		var errEx error
		value, errEx = conn.Do("GET", key)
		if errEx != nil {
			return errEx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return value, nil
}

// Expire sets key with expiration time
func (c *client) Expire(ctx context.Context, key string, expirationTimeout time.Duration) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		_, err := conn.Do("EXPIRE", key, int(expirationTimeout.Seconds()))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// HashSet sets hashset with name key
func (c *client) HashSet(ctx context.Context, key string, values interface{}) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		_, err := conn.Do("HSET", redis.Args{key}.AddFlat(values)...)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// HGetAll gets content of hashmap
func (c *client) HGetAll(ctx context.Context, key string) ([]interface{}, error) {
	var values []interface{}
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		var errEx error
		values, errEx = redis.Values(conn.Do("HGETALL", key))
		if errEx != nil {
			return errEx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return values, nil
}
