package simulation

import (
	"gopkg.in/ini.v1"
)

// Config is a set of configuration values for a simulation
type Config struct {
	Name           string
	Attacker       *AttackerConfig
	MainHandWeapon *WeaponConfig
	Feats          *FeatsConfig
	Target         *TargetConfig
}

// AttackerConfig is a configuration of high-level character attributes like BAB and Feats
type AttackerConfig struct {
	BaseAttackBonus     int
	StrMod              int
	DexMod              int
	AttacksPerRound     int
	BaseAttackBonusStep int
}

// WeaponConfig is a configuration for wieldable weapons by a character
type WeaponConfig struct {
	CritRange         int
	CritMultiplier    int
	Finessable        bool
	Ranged            bool
	AttackBonus       int
	BaseDamage        string
	DamageBonus       string
	PermEssenceDamage string
	TempEssenceDamage string
	AdditionalDamage1 string
	AdditionalDamage2 string
	ApplyTwoHandBonus bool
	Keen              bool
}

// TargetConfig is a configuration for what the attacker is hitting
type TargetConfig struct {
	ArmorClass int
}

const (
	//FeatImprovedCritical extends the critical range
	FeatImprovedCritical = "improved_critical"
	//FeatKiCritical adds 2 to the critical range
	FeatKiCritical = "ki_critical"
	//FeatIncreasedMultiplier adds 1 to the critical multiplier
	FeatIncreasedMultiplier = "increased_multiplier"
	//FeatRapidShot adds one attack at full attack bonus but subtracts 2 from every attack that round
	FeatRapidShot = "rapid_shot"
)

// FeatsConfig is a wrapper around a map with convenience functions
type FeatsConfig struct {
	all map[string]bool
}

// HasFeat tests to see if the Feat was taken
func (feats *FeatsConfig) HasFeat(name string) bool {
	value, ok := feats.all[name]
	if ok {
		return value
	}

	return false
}

// ConfigFromIni reads an ini file and creates a simulation.Config
func ConfigFromIni(file *ini.File) (cfg *Config, err error) {
	attackerSect := file.Section("attacker")
	mhSect := file.Section("main_hand_weapon")
	targetSect := file.Section("target")

	cfg = &Config{
		Name: file.Section("").Key("name").MustString("No name given"),
		Attacker: &AttackerConfig{
			BaseAttackBonus:     attackerSect.Key("base_attack_bonus").MustInt(0),
			StrMod:              attackerSect.Key("str_mod").MustInt(0),
			DexMod:              attackerSect.Key("dex_mod").MustInt(0),
			AttacksPerRound:     attackerSect.Key("attacks_per_round").MustInt(0),
			BaseAttackBonusStep: attackerSect.Key("base_attack_bonus_step").MustInt(-5),
		},
		MainHandWeapon: &WeaponConfig{
			CritRange:         mhSect.Key("crit_range").MustInt(1),
			CritMultiplier:    mhSect.Key("crit_multiplier").MustInt(2),
			AttackBonus:       mhSect.Key("attack_bonus").MustInt(0),
			BaseDamage:        mhSect.Key("base_damage").MustString("1d8"),
			DamageBonus:       mhSect.Key("damage_bonus").MustString("0"),
			PermEssenceDamage: mhSect.Key("perm_essence_damage").MustString("0"),
			TempEssenceDamage: mhSect.Key("temp_essence_damage").MustString("0"),
			Keen:              mhSect.Key("keen").MustBool(false),
			Finessable:        mhSect.Key("finessable").MustBool(false),
			ApplyTwoHandBonus: mhSect.Key("apply_two_hand_bonus").MustBool(false),
			AdditionalDamage1: mhSect.Key("additional_damage.1").MustString("0"),
			AdditionalDamage2: mhSect.Key("additional_damage.2").MustString("0"),
			Ranged:            mhSect.Key("ranged").MustBool(true),
		},
		Feats: &FeatsConfig{
			all: make(map[string]bool),
		},
		Target: &TargetConfig{
			ArmorClass: targetSect.Key("armor_class").MustInt(0),
		},
	}

	for _, featKey := range file.Section("feats").Keys() {
		cfg.Feats.all[featKey.Name()] = featKey.MustBool(false)
	}

	return cfg, nil
}
