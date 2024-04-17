package redis

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/go-redis/redis_rate/v10"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"

	goRedis "github.com/redis/go-redis/v9"
)

type Redis struct {
	redisRing       *goRedis.Ring
	rateLimiter     *redis_rate.Limiter
	redisCache      *cache.Cache
	redSync         *redsync.Redsync
	redisRingOnce   sync.Once
	rateLimiterOnce sync.Once
	redisCacheOnce  sync.Once
	redSyncOnce     sync.Once
}

type IRedis interface {
	GetRing() *goRedis.Ring
	SetLimiter()
	GetLimiter() *redis_rate.Limiter
	SetCache()
	GetCache() *cache.Cache
	SetRedSync()
	GetRedSync() *redsync.Redsync
}

var _ IRedis = (*Redis)(nil)

func NewRedis(opt *goRedis.RingOptions) IRedis {
	rt := &Redis{}
	rt.redisRingOnce.Do(func() {
		rt.redisRing = goRedis.NewRing(opt)
		_ = rt.redisRing.ForEachShard(context.TODO(),
			func(ctx context.Context, shard *goRedis.Client) error {
				// TODO: implement redisotel if metric infra ready
				//return redisotel.InstrumentMetrics(shard)
				return nil
			})
	})

	return rt
}

func (r *Redis) SetLimiter() {
	r.rateLimiterOnce.Do(func() {
		r.rateLimiter = redis_rate.NewLimiter(r.redisRing)
	})
}

func (r *Redis) SetCache() {
	r.redisCacheOnce.Do(func() {
		r.redisCache = cache.New(&cache.Options{
			Redis:      r.redisRing,
			LocalCache: cache.NewTinyLFU(1000, 5*time.Minute),
		})
	})
}

func (r *Redis) SetRedSync() {
	r.redSyncOnce.Do(func() {
		r.redSync = redsync.New(goredis.NewPool(r.redisRing))
	})
}

func (r *Redis) GetRing() *goRedis.Ring {
	return r.redisRing
}

func (r *Redis) GetLimiter() *redis_rate.Limiter {
	return r.rateLimiter
}

func (r *Redis) GetCache() *cache.Cache {
	return r.redisCache
}

func (r *Redis) GetRedSync() *redsync.Redsync {
	return r.redSync
}
