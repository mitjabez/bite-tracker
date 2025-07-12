package models

type MealView struct {
	MealType       string
	DateOfMeal     string
	TimeOfMeal     string
	Description    string
	HungerLevel    int32
	UsedSymptoms   []MealSymptom
	UnusedSymptoms []MealSymptom
}

type MealSymptom struct {
	Name  string
	Emoji string
}

var AllSymptoms = map[string]MealSymptom{
	"Bloating": {
		Name:  "Bloating",
		Emoji: "🎈",
	},
	"Gas": {
		Name:  "Gas",
		Emoji: "💨",
	},
	"Acid": {
		Name:  "Acid",
		Emoji: "🔥",
	},
	"Full": {
		Name:  "Full",
		Emoji: "🍽️",
	},
}
