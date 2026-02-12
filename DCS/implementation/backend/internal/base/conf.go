package base

import "time"

func GetTransactionTimeout() time.Duration {
	return 5 * time.Minute
}
