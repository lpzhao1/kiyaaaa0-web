package webv4

import (
	"fmt"
	"time"
)

type FilterBuilder func(next Filter) Filter

type Filter func(c *Context)

func MetricFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		startTime := time.Now().UnixNano()
		next(c)
		endTime := time.Now().UnixNano()
		fmt.Printf("run time: %d \n", endTime-startTime)
	}
}
