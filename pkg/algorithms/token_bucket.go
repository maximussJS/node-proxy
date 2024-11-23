package algorithms

import (
	"context"
	"fmt"
	"time"
)

type tokens chan struct{}

type TokenBucket struct {
	count  int64
	tokens tokens
	ticker *time.Ticker
}

func NewTokenBucket(count int64, rate int64) *TokenBucket {
	tokens := make(chan struct{}, count)

	c := int(count)
	for i := 0; i < c; i++ {
		tokens <- struct{}{}
	}

	fmt.Printf("Initialized TokenBucket with count %v and token rateMs %v\n", len(tokens), rate)
	everyMs := 1 / float64(rate) * 1000
	return &TokenBucket{
		count:  count,
		tokens: tokens,
		ticker: time.NewTicker(time.Duration(int64(everyMs) * int64(time.Millisecond))),
	}
}

func (tb *TokenBucket) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-tb.ticker.C:
				select {
				case tb.tokens <- struct{}{}:
				default:
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (tb *TokenBucket) Wait() {
	<-tb.tokens
}
