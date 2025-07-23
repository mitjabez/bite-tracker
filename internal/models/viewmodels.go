package models

import (
	"slices"
)

type MealView struct {
	Id             string
	MealType       string
	DateOfMeal     string
	TimeOfMeal     string
	Description    string
	HungerLevel    int32
	UsedSymptoms   []string
	UnusedSymptoms []string
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

func (mv MealView) IsSymptomUsed(symptom string) bool {
	return slices.Contains(mv.UsedSymptoms, symptom)
}
