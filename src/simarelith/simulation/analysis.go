package simulation

// Analysis is a high-level collection of interesting attributes
type Analysis struct {
	TotalAttacks       int
	FirstHitPercentage float32
	HitPercentage      float32
	CritPercentage     float32
	TotalDamage        int
	DamagePerRound     float32
}

// NewAnalysis returns a high-level analysis of combat rounds
func NewAnalysis(rounds []CombatRoundSimulationResult) Analysis {
	analysis := Analysis{}
	hits := 0
	firstHits := 0
	crits := 0
	for i := 0; i < len(rounds); i++ {
		analysis.TotalAttacks += len(rounds[i].attacks)
		if rounds[i].attacks[0].hit {
			firstHits++
		}
		hits += rounds[i].hits
		crits += rounds[i].crits
		analysis.TotalDamage += rounds[i].damage
	}

	analysis.HitPercentage = float32(hits) / float32(analysis.TotalAttacks)
	analysis.FirstHitPercentage = float32(firstHits) / float32(len(rounds))
	analysis.CritPercentage = float32(crits) / float32(analysis.TotalAttacks)
	analysis.DamagePerRound = float32(analysis.TotalDamage) / float32(len(rounds))

	return analysis
}
