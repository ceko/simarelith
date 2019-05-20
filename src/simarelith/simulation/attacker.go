package simulation

// Attacker has convenience methods for dealing with attack simulation data
type Attacker struct {
	config *Config
}

// AttackResult is the result of a single attack
type AttackResult struct {
	BaseAttackBonus    int
	Roll               int
	AttackBonusPenalty int
	AttributeMod       int // AttributeMod is adding dex to the roll if finessable or str if not
	WeaponAttackBonus  int
	TwoHandAttackBonus int
}

// AttacksPerRound returns how many attacks the attacker will perform every round, not influenced by feats
// or activated abilities (like rapid shot or flurry of blows)
func (attacker *Attacker) AttacksPerRound() int {
	return attacker.config.Attacker.AttacksPerRound
}

// ExtraAttacks is how many extra attacks are granted to the character through things like feats
func (attacker *Attacker) ExtraAttacks() int {
	extraAttacks := 0
	if attacker.config.Feats.HasFeat(FeatRapidShot) {
		extraAttacks++
	}

	return extraAttacks
}

// BaseAttackPenaltySchedule is an array of penalties to apply to attack rounds
func (attacker *Attacker) BaseAttackPenaltySchedule() []int {
	schedule := make([]int, attacker.AttacksPerRound())
	globalPenalty := 0
	if attacker.config.Feats.HasFeat(FeatRapidShot) {
		globalPenalty = -2
	}
	for i := 0; i < len(schedule); i++ {
		schedule[i] = i*attacker.config.Attacker.BaseAttackBonusStep + globalPenalty
	}

	return schedule
}

// ExtraAttackPenaltySchedule is an array of penalties to apply to attack rounds
func (attacker *Attacker) ExtraAttackPenaltySchedule() []int {
	schedule := make([]int, attacker.ExtraAttacks())
	cnt := 0
	if attacker.config.Feats.HasFeat(FeatRapidShot) {
		schedule[cnt] = -2
	}

	return schedule
}

// AttackPenaltySchedule is a combination of BaseAttackPenaltySchedule and ExtraAttackPenaltySchedule
func (attacker *Attacker) AttackPenaltySchedule() []int {
	return append(attacker.BaseAttackPenaltySchedule(), attacker.ExtraAttackPenaltySchedule()...)
}

// RollAttack gets the result of a single roll at penalty
func (attacker *Attacker) RollAttack(penalty int) AttackResult {
	result := AttackResult{
		BaseAttackBonus:    attacker.config.Attacker.BaseAttackBonus,
		AttackBonusPenalty: penalty,
	}

	result.Roll = roll("1d20")
	result.AttributeMod = attacker.config.Attacker.StrMod
	// If the weapon is finessable and not wielded as a 2h weapon, use dex instead
	if attacker.config.MainHandWeapon.Ranged ||
		(attacker.config.MainHandWeapon.Finessable && !attacker.config.MainHandWeapon.ApplyTwoHandBonus) {
		result.AttributeMod = attacker.config.Attacker.DexMod
	}
	result.WeaponAttackBonus = attacker.config.MainHandWeapon.AttackBonus
	// Can't get 2 hand attack bonus on bows
	if !attacker.config.MainHandWeapon.Ranged &&
		(!attacker.config.MainHandWeapon.Finessable && attacker.config.MainHandWeapon.ApplyTwoHandBonus) {
		// 2 AB for attacking with a two-hand weapon
		result.TwoHandAttackBonus = 2
	}

	return result
}

// HitsWith sees if an attack result is successful
func (attacker *Attacker) HitsWith(res *AttackResult, ac int) bool {
	// 1 always misses on a hit check
	if res.Roll == 1 {
		return false
	}
	// 20 always hits on a hit check
	if res.Roll == 20 {
		return true
	}

	return res.Total() >= ac
}

// CritConfirmedWith checks if a crit confirmation roll hits
func (attacker *Attacker) CritConfirmedWith(res *AttackResult, ac int) bool {
	// 1 is not a failure and 20 is not an automatic success
	return res.Total() >= ac
}

// CritsWith sees if an attack result is within critical threat
func (attacker *Attacker) CritsWith(res *AttackResult, threatRange int) bool {
	return res.Roll >= threatRange
}

// Total sums up the attack for the final attack roll result
func (res *AttackResult) Total() int {
	return res.BaseAttackBonus + res.Roll + res.AttackBonusPenalty + res.AttributeMod + res.WeaponAttackBonus + res.TwoHandAttackBonus
}
