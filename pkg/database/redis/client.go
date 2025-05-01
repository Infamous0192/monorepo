package redis

import (
	"app/pkg/database"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
	cfg    *database.DatabaseConfig
}

func NewClient(cfg database.DatabaseConfig) *Client {
	return &Client{cfg: &cfg}
}

func (c *Client) Connect() error {
	c.client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("redis://%s:%s@%s:%s", c.cfg.Username, c.cfg.Password, c.cfg.Host, c.cfg.Port),
		DB:   0,
	})
	return nil
}

func (c *Client) Disconnect() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func (c *Client) Ping() error {
	if c.client != nil {
		return c.client.Ping(context.Background()).Err()
	}
	return fmt.Errorf("redis connection not established")
}

// Get retrieves a value from Redis by key
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("redis connection not established")
	}
	return c.client.Get(ctx, key).Result()
}

// Set stores a value in Redis with an expiration
func (c *Client) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("redis connection not established")
	}
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Del removes one or more keys from Redis
func (c *Client) Del(ctx context.Context, keys ...string) error {
	if c.client == nil {
		return fmt.Errorf("redis connection not established")
	}
	return c.client.Del(ctx, keys...).Err()
}

// SAdd adds one or more members to a Redis set
func (c *Client) SAdd(ctx context.Context, key string, members ...interface{}) error {
	if c.client == nil {
		return fmt.Errorf("redis connection not established")
	}
	return c.client.SAdd(ctx, key, members...).Err()
}

// SRem removes one or more members from a Redis set
func (c *Client) SRem(ctx context.Context, key string, members ...interface{}) error {
	if c.client == nil {
		return fmt.Errorf("redis connection not established")
	}
	return c.client.SRem(ctx, key, members...).Err()
}

// SMembers returns all members of a Redis set
func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	if c.client == nil {
		return nil, fmt.Errorf("redis connection not established")
	}
	return c.client.SMembers(ctx, key).Result()
}

// SIsMember checks if a value is a member of a Redis set
func (c *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	if c.client == nil {
		return false, fmt.Errorf("redis connection not established")
	}
	return c.client.SIsMember(ctx, key, member).Result()
}

// Publish sends a message to a Redis channel
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	if c.client == nil {
		return fmt.Errorf("redis connection not established")
	}
	return c.client.Publish(ctx, channel, message).Err()
}

// Subscribe returns a Redis subscription for the given channels
func (c *Client) Subscribe(ctx context.Context, channels ...string) (*redis.PubSub, error) {
	if c.client == nil {
		return nil, fmt.Errorf("redis connection not established")
	}
	return c.client.Subscribe(ctx, channels...), nil
}

// PSubscribe returns a Redis subscription for the given patterns
func (c *Client) PSubscribe(ctx context.Context, patterns ...string) (*redis.PubSub, error) {
	if c.client == nil {
		return nil, fmt.Errorf("redis connection not established")
	}
	return c.client.PSubscribe(ctx, patterns...), nil
}

// ScanKeys returns all keys matching a pattern
func (c *Client) ScanKeys(ctx context.Context, pattern string) ([]string, error) {
	if c.client == nil {
		return nil, fmt.Errorf("redis connection not established")
	}

	var keys []string
	var cursor uint64
	var err error

	for {
		var scanKeys []string
		scanKeys, cursor, err = c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}
