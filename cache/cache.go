package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	redis "github.com/go-redis/redis/v8"
	goLogging "github.com/op/go-logging"
)

// Context is the runtime context for access Redis
type Context struct {
	URL        string
	Logger     *goLogging.Logger
	Redis      *redis.Client
	MaxRetries int
	Wait       int
	Locker     *redislock.Client
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

	ctx.Locker = redislock.New(ctx.Redis)

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
			ctx.Logger.Debugf("[%s] Checking if entry has a cache lock, attempt #%d", id, i+1)
			if ret, err := ctx.Redis.HGet(lckCtx, id, "lock").Result(); err != redis.Nil {
				expiry, err := time.Parse(time.RFC3339, ret)
				if err != nil {
					ctx.Logger.Errorf("[%s] Could not parse expiry of cache entry %s: %s", id, expiry, err.Error())
					ctx.Clear(lckCtx, id)
					break
				}
				if expiry.Before(time.Now()) {
					ctx.Logger.Errorf("[%s] The lock for this entry has expired")
					ctx.Clear(lckCtx, id)
					break
				}
				ctx.Logger.Warningf("[%s] a lock was found in the cache for document, sleeping for %d seconds", id, Wait)
				time.Sleep(time.Second * Wait)
			} else {
				break
			}
		}
	}

	ctx.Logger.Debugf("[%s] Retrieving cache entry", id)
	cached, err := ctx.Redis.HGet(lckCtx, id, "data").Result()
	if err == redis.Nil {
		ctx.Logger.Debugf("[%s] Entry not found in cache", id)
		return "", errors.New("not found")
	}

	ctx.Logger.Debugf("[%s] Entry found in cache", id)

	return cached, nil
}

// Save commits a key/value pair into Redis
func (ctx *Context) Save(lckCtx context.Context, id string, data interface{}) error {

	ctx.Logger.Debugf("[%s] Saving cache entry", id)
	str, err := json.Marshal(data)
	if err != nil {
		ctx.Logger.Errorf("[%s] failed to serialise data: %s", id, err.Error())
		return fmt.Errorf("failed to serialise data: %s", err.Error())
	}

	ctx.Logger.Debugf("[%s] Writing to Redis", id)
	if _, err = ctx.Redis.HSet(lckCtx, id, "data", str).Result(); err != nil {
		ctx.Logger.Errorf("[%s] Failed to write to Redis: %s", id, err.Error())
		return fmt.Errorf("failed to save document in cache: %s", err.Error())
	}
	ctx.Logger.Debugf("[%s] Cache write was successful", id)
	return nil
}

// Lock puts a special marker into a redis hash entry preventing future access
func (ctx *Context) Lock(lckCtx context.Context, id string) (*redislock.Lock, error) {

	ctx.Logger.Debugf("[%s] Locking cache entry", id)
	retryStrategy := redislock.LimitRetry(redislock.ExponentialBackoff(LckRetryTTLMin*time.Millisecond, LckRetryTTLMax*time.Millisecond), LckRetryCount)
	// Try to obtain lock for the connection for 300 ms
	return ctx.Locker.Obtain(lckCtx, id, LckLockTTL*time.Second, &redislock.Options{RetryStrategy: retryStrategy})
}

// Clear removes a Redis cache entry identified by the "id" parameter
func (ctx *Context) Clear(lckCtx context.Context, id string) error {

	ctx.Logger.Debugf("[%s] Removing cache entry", id)
	if _, err := ctx.Redis.Del(lckCtx, id).Result(); err != nil {
		ctx.Logger.Errorf("[%s] Failed to remove cache entry: %s", id, err.Error())
		return fmt.Errorf("failed to remove document in cache: %s", err.Error())
	}
	ctx.Logger.Debugf("[%s] Successfully removed the record", id)
	return nil
}
