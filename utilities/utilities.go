package utilities

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"
)

func Retry(
	ctx context.Context,
	times int,
	baseDelay time.Duration,
	fn func(ctx context.Context) error,
) error {
	var err error

	for counter := 0; counter < times; counter++ {
		// Stop immediately if context is cancelled/deadline exceeded
		if ctx.Err() != nil {
			return ctx.Err()
		}

		backoffMultiplier := math.Exp2(float64(counter))

		currDelay := time.Duration(float64(baseDelay) * backoffMultiplier)

		if err = fn(ctx); err != nil {
			log.Printf("[Retry] Attempt %v, backing off for %v, failed with error: %v.",
				counter+1, currDelay, err)

			timer := time.NewTimer(currDelay)

			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()

				return ctx.Err()
			}

			continue
		}

		break
	}

	if err != nil {
		return fmt.Errorf("retry attempts exhausted after %d tires: %w", times, err)
	}

	return nil
}
