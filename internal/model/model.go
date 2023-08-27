package model

import "time"

type Segment struct {
	SegmentName string
	CreatedAt   time.Time
}

type User struct {
	UserID    int
	Username  string
	CreatedAt time.Time
}
