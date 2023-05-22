package params

import "time"

type TimestampMicro int64

func (t *TimestampMicro) GetTimestamp() int64 {
	return int64(*t) / 1000
}

func (t *TimestampMicro) GetLocalTime() time.Time {
	timeInt := int64(*t)
	return time.Unix(timeInt/1000, 0).Local()
}
