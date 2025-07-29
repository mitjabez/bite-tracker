package model

import (
	"slices"
	"time"
)

type MealView struct {
	Id          string
	MealType    string
	TimeOfMeal  time.Time
	Description string
	HungerLevel int32
	Symptoms    []string
}

type MealSymptom struct {
	Name  string
	Emoji string
}

var Symptoms = []string{
	"acid",
	"bloating",
	"full",
	"gas",
}

// TODO: Optimize if needed
func (mv MealView) IsSymptomUsed(symptom string) bool {
	return slices.Contains(mv.Symptoms, symptom)
}

func ResolveMealType(time time.Time) string {
	hour := time.Hour()
	switch {
	case hour < 9:
		return "breakfast"
	case hour < 11:
		return "brunch"
	case hour < 15:
		return "lunch"
	default:
		return "dinner"
	}
}
