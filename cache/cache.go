package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/9spokes/go/logging/v3"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	redis "github.com/go-redis/redis/v8"
)

// Context is the runtime context for access Redis
type Context struct {
	URL        string
	Logger     *logging.Logger
	Redis      *redis.Client
	MaxRetries int
	Wait       int
	RedSync    *redsync.Redsync
}

// MaxRetries is the number of times we re-attempt to access the cache when it is locked
const MaxRetries = 20

// Wait is the amount of time (in seconds) we wait before trying
const Wait = 2

const (
	LckRetryTTLMin = 50  //MS
	LckRetryTTLMax = 600 //MS
	LckRetryCount  = 200
	LckLockTTL     = 10 //Sec
)

// New creates a new instance of a Redis cache and returns a context for future use
func New(url string) (*Context, error) {
	redisOpts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %s", err.Error())
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisOpts.Addr,
		Password:     redisOpts.Password,
		DB:           redisOpts.DB,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     100,
		PoolTimeout:  30 * time.Second,
	})

	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	ctx := Context{
		Redis:      client,
		URL:        url,
		MaxRetries: MaxRetries,
		Wait:       Wait,
		RedSync:    rs,
	}

	_, err = ctx.Redis.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %s", err.Error())
	}

	return &ctx, nil
}

// Get grabs an entry from the Redis cache matching the key identified by the "id" parameter and returns the associated
// unmarkshaled document. If lock is true it first checks if there is a lock on the entry and if found waits until the
// lock is released.
func (ctx *Context) Get(lckCtx context.Context, id string, lock bool) (string, error) {

	if lock {
		for i := 0; i < ctx.MaxRetries; i++ {
			logging.Debugf("[%s] Checking if entry has a cache lock, attempt #%d", id, i+1)
			if ret, err := ctx.Redis.HGet(lckCtx, id, "lock").Result(); err != redis.Nil {
				expiry, err := time.Parse(time.RFC3339, ret)
				if err != nil {
					logging.Errorf("[%s] Could not parse expiry of cache entry %s: %s", id, expiry, err.Error())
					ctx.Clear(lckCtx, id)
					break
				}
				if expiry.Before(time.Now()) {
					logging.Errorf("[%s] The lock for this entry has expired")
					ctx.Clear(lckCtx, id)
					break
				}
				logging.Warningf("[%s] a lock was found in the cache for document, sleeping for %d seconds", id, Wait)
				time.Sleep(time.Second * Wait)
			} else {
				break
			}
		}
	}

	logging.Debugf("[%s] Retrieving cache entry", id)
	cached, err := ctx.Redis.HGet(lckCtx, id, "data").Result()
	if err == redis.Nil {
		logging.Debugf("[%s] Entry not found in cache", id)
		return "", errors.New("not found")
	}

	logging.Debugf("[%s] Entry found in cache", id)

	return cached, nil
}

// Save commits a key/value pair into Redis
func (ctx *Context) Save(lckCtx context.Context, id string, data interface{}) error {

	logging.Debugf("[%s] Saving cache entry", id)
	str, err := json.Marshal(data)
	if err != nil {
		logging.Errorf("[%s] failed to serialise data: %s", id, err.Error())
		return fmt.Errorf("failed to serialise data: %s", err.Error())
	}

	logging.Debugf("[%s] Writing to Redis", id)
	if _, err = ctx.Redis.HSet(lckCtx, id, "data", str).Result(); err != nil {
		logging.Errorf("[%s] Failed to write to Redis: %s", id, err.Error())
		return fmt.Errorf("failed to save document in cache: %s", err.Error())
	}
	logging.Debugf("[%s] Cache write was successful", id)
	return nil
}

func (ctx *Context) Lock(id string) (func(), error) {

	mutex := ctx.RedSync.NewMutex(id)

	for i := 0; i < ctx.MaxRetries; i++ {
		err := mutex.Lock()
		//Failed to acquire lock after exhausting all retries, keep try until its unlocked
		if err == redsync.ErrFailed {
			logging.Debugf("failed to acquire lock after exhausting all retries, sleeping for %d seconds before retrying", Wait)
			time.Sleep(time.Second * Wait)
			continue
		}
		if err != nil {
			return nil, err
		}
		break
	}

	return func() {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			logging.Errorf("failed to unlock, %s", err.Error())
		}
	}, nil
}

// Clear removes a Redis cache entry identified by the "id" parameter
func (ctx *Context) Clear(lckCtx context.Context, id string) error {

	logging.Debugf("[%s] Removing cache entry", id)
	if _, err := ctx.Redis.Del(lckCtx, id).Result(); err != nil {
		logging.Errorf("[%s] Failed to remove cache entry: %s", id, err.Error())
		return fmt.Errorf("failed to remove document in cache: %s", err.Error())
	}
	logging.Debugf("[%s] Successfully removed the record", id)
	return nil
}
