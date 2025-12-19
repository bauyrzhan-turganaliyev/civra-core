package entity

type Profession string
type Resource string

const (
	ProfFarmer  Profession = "farmer"
	ProfMiner   Profession = "miner"
	ProfLumber  Profession = "lumber"
	ProfKing    Profession = "king"
	ProfBuilder Profession = "builder"
	ProfSmith   Profession = "blacksmith"

	ResFood Resource = "food"
	ResIron Resource = "iron"
	ResWood Resource = "wood"
)
