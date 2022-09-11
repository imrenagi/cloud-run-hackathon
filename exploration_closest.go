package main

func NewClosestDistanceExploration(p *Player) *Exploration {
	return &Exploration{
		p: p,
		sorter: StepSorter{p: *p},
	}
}

func NewWeightedExploration(p *Player) *Exploration {
	return &Exploration{
		p: p,
		sorter: WeightedSorter{p: *p},
	}
}


