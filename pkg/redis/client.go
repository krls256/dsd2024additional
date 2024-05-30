package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"time"
)

var Nil = redis.Nil

type Client struct {
	cfg   *Config
	redis *redis.Client
}

func New(cfg *Config) (client *Client, err error) {
	client = &Client{
		cfg: cfg,
		redis: redis.NewClient(&redis.Options{
			Addr:     cfg.Host + ":" + cfg.Port,
			Password: cfg.Pass,
		})}

	if _, err = client.redis.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.redis.Set(ctx, c.prepareKey(key), value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	return c.redis.Get(ctx, c.prepareKey(key)).Bytes()
}

// FIFO queue start

func (c *Client) LLen(ctx context.Context, key string) (int64, error) {
	return c.redis.LLen(ctx, c.prepareKey(key)).Result()
}

func (c *Client) LPush(ctx context.Context, key, value string) error {
	return c.redis.LPush(ctx, c.prepareKey(key), value).Err()
}

func (c *Client) RPop(ctx context.Context, key string) (string, error) {
	return c.redis.RPop(ctx, c.prepareKey(key)).Result()
}

// FIFO queue end

func (c *Client) Del(ctx context.Context, keys ...string) error {
	for i, v := range keys {
		keys[i] = c.prepareKey(v)
	}

	return c.redis.Del(ctx, keys...).Err()
}

func (c *Client) prepareKey(s string) string {
	return fmt.Sprintf("%s:%s", c.cfg.Prefix, s)
}

func (c *Client) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.redis.Subscribe(ctx, lo.Map(channels, func(item string, index int) string {
		return c.prepareKey(item)
	})...)
}

func (c *Client) Publish(ctx context.Context, channel string, message any) error {
	_, err := c.redis.Publish(ctx, c.prepareKey(channel), message).Result()

	return err
}
