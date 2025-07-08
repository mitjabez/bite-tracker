package models

import (
	"time"
)

type Meal struct {
	Id             string
	Type           string
	Time           time.Time
	Description    string
	HungerLevel    int64
	Symptoms       string
	FeltFineAfter  bool
	HadIssuesAfter bool
}
