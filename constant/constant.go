package constant

import "fmt"

type redisK string

func (r redisK) Format(data ...interface{}) string {
	return fmt.Sprintf(string(r), data...)
}

const (
	SpikeKey redisK = "stock_%d"

	EtcdKey = "Spike_Server"
)
