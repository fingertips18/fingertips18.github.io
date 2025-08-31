package domain

import "time"

type TimeProvider func() time.Time
