package rate

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// 令牌桶：桶、填充速率、消费

type Rate struct {
	capacity  float64 // 桶容量
	tokens    float64 // 桶内令牌数量
	rate      float64 // 填充速率,每 time.Duration需要填充多少个令牌
	rw        sync.Mutex
	ctx       context.Context
	lastTime  time.Time     //最后填充时间
	precision time.Duration // 填充精度，N分之一秒填充一个令牌
	timeout   time.Duration //wait等待令牌超时时间
}

type Option func(opts *Rate)

const (
	defaultCapacity  = 1000000
	defaultMax       = 1000000
	defaultSecond    = time.Second
	defaultTimeout   = time.Second
	defaultPrecision = time.Millisecond * 100
)

var (
	ErrRateLimit = errors.New("rate limit")
)

// New 每 second 秒内运行放行最大 max 个请求，桶的容量最大为 capacity
// capacity 桶容量
// max N秒内的允许最大请求数
// second N秒
func New(ctx context.Context, opts ...Option) *Rate {
	r := &Rate{
		rw:       sync.Mutex{},
		ctx:      ctx,
		lastTime: time.Now(),
	}

	defaultOpts := []Option{
		WithCapacity(defaultCapacity),
		WithRate(defaultMax, defaultSecond),
		WithTimeout(defaultTimeout),
		WithPrecision(defaultPrecision),
	}

	opts = append(defaultOpts, opts...)

	for _, opt := range opts {
		opt(r)
	}

	fmt.Println(r)

	go r.refill()
	return r
}

func WithCapacity(capacity int64) Option {
	return func(opts *Rate) {
		opts.capacity = float64(capacity)
		opts.tokens = opts.capacity
	}
}

func WithRate(max, second time.Duration) Option {
	return func(opts *Rate) {
		opts.rate = float64(max) / float64(second*time.Second)
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opts *Rate) {
		opts.timeout = timeout
	}
}

func WithPrecision(precision time.Duration) Option {
	return func(opts *Rate) {
		opts.precision = precision
	}
}

func (r *Rate) refill() {
	ticker := time.NewTicker(r.precision)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.fill(time.Now())
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *Rate) fill(now time.Time) {
	r.rw.Lock()
	defer r.rw.Unlock()

	sub := now.Sub(r.lastTime)
	if sub <= 0 {
		return
	}

	r.lastTime = now
	r.tokens += float64(sub) * r.rate
	if r.tokens > r.capacity {
		r.tokens = r.capacity
	}
}

// Allow 是否允许放行
func (r *Rate) Allow() bool {
	return r.AllowN(time.Now(), 1)
}

func (r *Rate) AllowN(t time.Time, n int64) bool {
	select {
	case <-r.ctx.Done():
		return false
	default:
		r.fill(t)
		r.rw.Lock()
		defer r.rw.Unlock()
		v := float64(n)
		if r.tokens >= v {
			r.tokens -= v
			return true
		}
		return false
	}
}

func (r *Rate) WaitN(ctx context.Context, n int64) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(r.timeout):
			if r.AllowN(time.Now(), n) {
				return nil
			}
			return ErrRateLimit
		}
	}
}
