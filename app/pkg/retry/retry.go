package retry

import (
	"context"
	"math/rand"
	"time"

	"github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/backoff"
	"github.com/kamilsk/retry/v5/jitter"
	"github.com/kamilsk/retry/v5/strategy"
)

func Retry(
	ctx context.Context,
	totalRetries int,
	action func(ctx context.Context) error,
) error {
	how := retry.How{
		strategy.Limit(uint(totalRetries)),
		strategy.BackoffWithJitter(
			backoff.Fibonacci(10*time.Millisecond),
			jitter.NormalDistribution(
				rand.New(rand.NewSource(time.Now().UnixNano())),
				0.25,
			),
		),
	}
	return retry.Do(ctx, action, how...)
}
