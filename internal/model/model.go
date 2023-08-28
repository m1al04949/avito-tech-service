package model

import "time"

type Segments struct {
	SegmentName string
	CreatedAt   time.Time
}

type Users struct {
	UserID    int
	CreatedAt time.Time
}

type UserSegments struct {
	UserID      int
	SegmentName string
}

type Segment struct {
	Slug string `json:"slug,omitempty"`
}
