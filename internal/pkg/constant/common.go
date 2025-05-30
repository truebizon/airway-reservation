package constant

import (
	"time"
)

const (
	REQUEST_TIMEOUT           = 1
	SLOW_THRESH_HOLD_QUERY    = time.Second
	TIME_FORMAT               = "2006-01-02 15:04:05"
	TIME_MICRO_FORMAT         = "2006/01/02 15:04:05.000000"
	TIME_MICROZ_FORMAT        = "2006/01/02 15:04:05.000000-07"
	TIME_NANO_FORMAT          = "2006/01/02 15:04:05.000000000"
	SIZE_DEFAULT              = 20
	PAGE_DEFAULT              = 1
	PAGENATION_SIZE_UNLIMITED = -1
	TIME_NANOZ_FORMAT_MST     = "2006-01-02 15:04:05.999999999 -0700 MST"
)
