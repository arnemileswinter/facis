package base

import "time"

func TransactionTimeout() time.Duration {
	return 5 * time.Minute
}
