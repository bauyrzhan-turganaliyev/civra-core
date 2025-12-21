package service

type ToolSpec struct {
	Tier     int
	BonusPct int
	MaxDur   int
	IronCost int
	WoodCost int
}

func ToolSpecForUI(tier int) (ToolSpec, bool) {
	return toolSpec(tier)
}

func toolSpec(tier int) (ToolSpec, bool) {
	switch tier {
	case 1:
		return ToolSpec{Tier: 1, BonusPct: 10, MaxDur: 50, IronCost: 2, WoodCost: 1}, true
	case 2:
		return ToolSpec{Tier: 2, BonusPct: 20, MaxDur: 100, IronCost: 4, WoodCost: 2}, true
	case 3:
		return ToolSpec{Tier: 3, BonusPct: 35, MaxDur: 150, IronCost: 8, WoodCost: 4}, true
	default:
		return ToolSpec{}, false
	}
}
