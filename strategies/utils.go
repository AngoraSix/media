package strategies

import (
	"fmt"
	"time"
)

func createNowString() string {
	now := time.Now()
	// we'll format it as yyyyMMddHHmmssSSS
	millis := fmt.Sprintf("%d", now.Nanosecond())
	millis = millis[:len(millis)-6]
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d%s", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), millis)
}
