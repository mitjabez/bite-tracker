package models

type MealSymptom string

const (
	Bloating MealSymptom = "bloating"
	Gas      MealSymptom = "gas"
	Acid     MealSymptom = "acid"
	Full     MealSymptom = "full"
)

var AllSymptoms = []MealSymptom{Bloating, Gas, Acid, Full}
