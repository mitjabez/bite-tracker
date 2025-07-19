package models

type MealView struct {
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
	"Acid",
	"Bloating",
	"Full",
	"Gas",
}
