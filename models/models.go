package models

import (
	"time"
)

type MealSymptom string

const (
	Bloating MealSymptom = "bloating"
	Gas      MealSymptom = "gas"
	Acid     MealSymptom = "acid"
	Full     MealSymptom = "full"
)

var AllSymptoms = []MealSymptom{Bloating, Gas, Acid, Full}

type Meal struct {
	Id          string
	Type        string
	Time        time.Time
	Description string
	HungerLevel int64
	Symptoms    []MealSymptom
}
