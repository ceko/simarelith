package simulation

// Weapon has convenience methods for dealing with weapon simulation data
type Weapon struct {
	config       *Config
	weaponConfig *WeaponConfig
}

// DamageRoll is a weapon damage roll
type DamageRoll struct {
	BaseDamage        int
	StrDamage         int
	EnhancementDamage int
}

// ModifiedCriticalThreat is the number you have to meet or beat to confirm a crit
func (weapon *Weapon) ModifiedCriticalThreat() int {
	// If the crit range is 1 this will make 20 a crit
	criticalThreat := 21 - weapon.weaponConfig.CritRange
	if weapon.config.Feats.HasFeat(FeatImprovedCritical) {
		criticalThreat -= weapon.weaponConfig.CritRange
	}
	// Ki critical adds a flat 2 to the crit range
	if weapon.config.Feats.HasFeat(FeatKiCritical) {
		criticalThreat -= 2
	}
	// Keen weapon works just like improved critical
	if weapon.weaponConfig.Keen {
		criticalThreat -= weapon.weaponConfig.CritRange
	}

	return criticalThreat
}

// ModifiedCritMultiplier returns how much the weapon should multiply its damage by when critting
func (weapon *Weapon) ModifiedCritMultiplier() int {
	multiplier := weapon.weaponConfig.CritMultiplier
	if weapon.config.Feats.HasFeat(FeatIncreasedMultiplier) {
		multiplier++
	}

	return multiplier
}

// RollDamage rolls and sums up all the different damage types
func (weapon *Weapon) RollDamage() DamageRoll {
	damage := DamageRoll{}
	damage.BaseDamage = roll(weapon.weaponConfig.BaseDamage)
	// no support for mighty
	if !weapon.weaponConfig.Ranged {
		damage.StrDamage = weapon.config.Attacker.StrMod
	}
	if !weapon.weaponConfig.Finessable && weapon.weaponConfig.ApplyTwoHandBonus {
		damage.StrDamage = int(float32(damage.StrDamage) * 1.5)
	}

	damage.EnhancementDamage = roll(weapon.weaponConfig.DamageBonus) +
		roll(weapon.weaponConfig.PermEssenceDamage) +
		roll(weapon.weaponConfig.TempEssenceDamage) +
		roll(weapon.weaponConfig.AdditionalDamage1) +
		roll(weapon.weaponConfig.AdditionalDamage2)

	return damage
}
