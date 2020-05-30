package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	goLogging "github.com/op/go-logging"
)

// Context is the runtime context for access Redis
type Context struct {
	URL        string
	Logger     *goLogging.Logger
	Redis      *redis.Client
	MaxRetries int
	Wait       int
}

// MaxRetries is the number of times we re-attempt to access the cache when it is locked
const MaxRetries = 20

// Wait is the amount of time (in seconds) we wait before trying
const Wait = 2

// New creates a new instance of a Redis cache and returns a context for future use
func New(url string, logger *goLogging.Logger) (*Context, error) {
	redisOpts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %s", err.Error())
	}

	if logger == nil {
		return nil, fmt.Errorf("a logger is required")
	}

	ctx := Context{
		Redis: redis.NewClient(&redis.Options{
			Addr:         redisOpts.Addr,
			Password:     redisOpts.Password,
			DB:           redisOpts.DB,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			PoolSize:     10,
			PoolTimeout:  30 * time.Second,
		}),
		URL:        url,
		Logger:     logger,
		MaxRetries: MaxRetries,
		Wait:       Wait,
	}

	_, err = ctx.Redis.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Redis: %s", err.Error())
	}

	return &ctx, nil
}

// Get grabs an entry from the Redis cache matching the key identified by the "id" parameter and returns the associated unmarkshaled document
func (ctx *Context) Get(id string) (string, error) {

	for i := 0; i < ctx.MaxRetries; i++ {
		ctx.Logger.Debugf("[%s] Checking if entry has a cache lock, attempt #%d", id, i+1)
		if ret, err := ctx.Redis.HGet(id, "lock").Result(); err != redis.Nil {
			expiry, err := time.Parse(time.RFC3339, ret)
			if err != nil {
				ctx.Logger.Errorf("[%s] Could not parse expiry of cache entry %s: %s", id, expiry, err.Error())
				ctx.Clear(id)
				break
			}
			if expiry.Before(time.Now()) {
				ctx.Logger.Errorf("[%s] The lock for this entry has expired")
				ctx.Clear(id)
				break
			}
			ctx.Logger.Warningf("[%s] a lock was found in the cache for document, sleeping for %d seconds", Wait, id)
			time.Sleep(time.Second * Wait)
		} else {
			break
		}
	}

	ctx.Logger.Debugf("[%s] Retrieving cache entry", id)
	cached, err := ctx.Redis.HGet(id, "data").Result()
	if err == redis.Nil {
		ctx.Logger.Debugf("[%s] Entry not found in cache", id)
		return "", errors.New("Not found")
	}

	ctx.Logger.Debugf("[%s] Entry found in cache", id)

	return cached, nil
}

// Save commits a key/value pair into Redis
func (ctx *Context) Save(id string, data interface{}) error {

	ctx.Logger.Debugf("[%s] Saving cache entry", id)
	str, err := json.Marshal(data)
	if err != nil {
		ctx.Logger.Errorf("[%s] failed to serialise data: %s", id, err.Error())
		return fmt.Errorf("failed to serialise data: %s", err.Error())
	}

	ctx.Logger.Debugf("[%s] Writing to Redis", id)
	if _, err = ctx.Redis.HSet(id, "data", str).Result(); err != nil {
		ctx.Logger.Errorf("[%s] Failed to write to Redis: %s", id, err.Error())
		return fmt.Errorf("failed to save document in cache: %s", err.Error())
	}
	ctx.Logger.Debugf("[%s] Cache write was successful", id)
	return nil
}

// Lock puts a special marker into a redis hash entry preventing future access
func (ctx *Context) Lock(id string) error {

	ctx.Logger.Debugf("[%s] Locking cache entry", id)
	if _, err := ctx.Redis.HSet(id, "lock", time.Now().Add(time.Duration(1*time.Minute)).Format(time.RFC3339)).Result(); err != nil {
		ctx.Logger.Debugf("[%s] Failed to lock cache entry: %s", id, err.Error())
		return fmt.Errorf("failed to lock document in cache: %s", err.Error())
	}
	ctx.Logger.Debugf("[%s] Successfully locked the record", id)
	return nil
}

// Unlock removes a special marker from a redis hash entry that was put there by Lock()
func (ctx *Context) Unlock(id string) error {

	ctx.Logger.Debugf("[%s] Unlocking cache entry", id)
	if _, err := ctx.Redis.HDel(id, "lock").Result(); err != nil {
		ctx.Logger.Errorf("[%s] Failed to unlock cache entry: %s", id, err.Error())
		return fmt.Errorf("failed to unlock document in cache: %s", err.Error())
	}
	ctx.Logger.Debugf("[%s] Successfully unlocked the record", id)
	return nil
}

// Clear removes a Redis cache entry identified by the "id" parameter
func (ctx *Context) Clear(id string) error {

	ctx.Logger.Debugf("[%s] Removing cache entry", id)
	if _, err := ctx.Redis.Del(id).Result(); err != nil {
		ctx.Logger.Errorf("[%s] Failed to remove cache entry: %s", id, err.Error())
		return fmt.Errorf("failed to remove document in cache: %s", err.Error())
	}
	ctx.Logger.Debugf("[%s] Successfully removed the record", id)
	return nil
}
