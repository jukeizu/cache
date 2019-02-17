package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
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
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.Db,
	})

	return &cache{redisClient, config.Version}
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

	redisError := c.client.Set(jsonKey, jsonValue, expiration)

	return redisError.Err()
}

func (c *cache) Get(key interface{}, value interface{}) error {
	jsonKey, err := c.getVersionedKey(key)
	if err != nil {
		return err
	}

	val, err := c.client.Get(jsonKey).Result()
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
