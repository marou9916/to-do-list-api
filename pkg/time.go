package pkg

import "time"

var TimeNow = func() time.Time {
	return time.Now()
}
