package models

import (
	"time"
)

type MealSymptom string

const (
	Bloating MealSymptom = "Bloating"
	Gas      MealSymptom = "Gas"
	Acid     MealSymptom = "Acid"
	Full     MealSymptom = "Full"
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
