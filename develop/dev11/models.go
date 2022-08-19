package main

import "time"

type Event struct {
	UserID      string
	Date        time.Time
	Description string
}
