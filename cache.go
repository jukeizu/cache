package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	DefaultRedisAddress = "localhost:6379"
)

type Config struct {
	Address  string
	Password string
	Db       int
	Version  string
}

type Cache interface {
	Set(interface{}, interface{}, time.Duration) error
	Get(interface{}, interface{}) error
}

type cache struct {
	client  *redis.Client
	version string
}

type versionedKey struct {
	Version string
	Content interface{}
}

func New(config Config) Cache {
	options := &redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.Db,
	}

	return NewWithRedisOptions(config.Version, options)
}

func NewWithRedisURL(version string, url string) (Cache, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return NewWithRedisOptions(version, options), nil
}

func NewWithRedisOptions(version string, options *redis.Options) Cache {
	return NewWithRedisClient(version, redis.NewClient(options))
}

func NewWithRedisClient(version string, client *redis.Client) Cache {
	return &cache{client, version}
}

func (c *cache) Set(key interface{}, value interface{}, expiration time.Duration) error {
	jsonKey, err := c.getVersionedKey(key)
	if err != nil {
		return err
	}

	jsonValue, err := asJsonString(value)
	if err != nil {
		return err
	}

	redisError := c.client.Set(context.Background(), jsonKey, jsonValue, expiration)

	return redisError.Err()
}

func (c *cache) Get(key interface{}, value interface{}) error {
	jsonKey, err := c.getVersionedKey(key)
	if err != nil {
		return err
	}

	val, err := c.client.Get(context.Background(), jsonKey).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), &value)
}

func (c *cache) getVersionedKey(key interface{}) (string, error) {
	v := versionedKey{
		Version: c.version,
		Content: key,
	}

	return asJsonString(v)
}

func asJsonString(val interface{}) (string, error) {
	marshalled, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(marshalled), nil
}
