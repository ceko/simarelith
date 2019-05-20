package simulation

import (
	"sync"

	"simarelith/logger"

	"github.com/davecgh/go-spew/spew"
)

// Simulator takes a config and runs simulations, returning results
type Simulator struct {
	config *Config
}

// NewSimulator returns a new Simulator from a config
func NewSimulator(cfg Config) Simulator {
	return Simulator{
		config: &cfg,
	}
}

// AttackSimulationResult is a collection of statistics for one attack
type AttackSimulationResult struct {
	round                    int
	attackRoll               int
	criticalThreatRange      int
	criticalThreat           bool
	modifiedAttackRoll       int
	hit                      bool
	criticalConfirmationRoll int
	criticalConfirmed        bool
	damage                   *DamageRoll
	damageMultiplier         int
	finalDamage              int
}

// CombatRoundSimulationResult is a collection of AttackSimulationResult
type CombatRoundSimulationResult struct {
	attacks []*AttackSimulationResult
	hits    int
	crits   int
	damage  int
}

// Run simulates combat and returns simulation results for aggregation
func (simulator *Simulator) Run(iterations int) []CombatRoundSimulationResult {
	logger.Trace.Println("Running", iterations, "simulations in parallel")
	rounds := make([]CombatRoundSimulationResult, iterations)
	wg := sync.WaitGroup{}
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(simulator *Simulator, iteration int, wg *sync.WaitGroup) {
			defer wg.Done()
			logger.Trace.Println("Iteration", iteration, "started")
			round := runCombatSimulation(simulator.config)
			//logger.Trace.Println("combat round: \n", spew.Sdump(round))
			rounds[iteration] = round
		}(simulator, i, &wg)
	}
	wg.Wait()
	logger.Trace.Println("Iterations complete")

	return rounds
}

func processAttack(round int, penalty int, config *Config, attacker *Attacker, mhWeapon *Weapon) *AttackSimulationResult {
	result := AttackSimulationResult{
		damageMultiplier: 1,
		round:            round,
	}

	// Roll to hit
	hitAttack := attacker.RollAttack(penalty)
	result.attackRoll = hitAttack.Roll
	result.modifiedAttackRoll = hitAttack.Total()

	// Check if attack hits against target AC
	if attacker.HitsWith(&hitAttack, config.Target.ArmorClass) {
		result.hit = true
		// See if the hit is a critical threat
		result.criticalThreat = attacker.CritsWith(&hitAttack, mhWeapon.ModifiedCriticalThreat())
		result.criticalThreatRange = mhWeapon.ModifiedCriticalThreat()
		damageRoll := mhWeapon.RollDamage()

		if result.criticalThreat {
			// Roll to confirm crit
			critConfirmationAttack := attacker.RollAttack(0)
			result.criticalConfirmationRoll = critConfirmationAttack.Total()
			if attacker.CritConfirmedWith(&critConfirmationAttack, config.Target.ArmorClass) {
				result.criticalConfirmed = true
				// Find out the damage multiplier, taking into account feats
				result.damageMultiplier = mhWeapon.ModifiedCritMultiplier()
			}
		}

		result.damage = &damageRoll
		result.finalDamage = (result.damage.BaseDamage + result.damage.StrDamage + result.damage.EnhancementDamage) * result.damageMultiplier
	} else {
		result.damage = &DamageRoll{}
	}

	logger.Trace.Println(spew.Sdump(result))
	return &result
}

func runCombatSimulation(config *Config) CombatRoundSimulationResult {
	attacker := &Attacker{config}
	mhWeapon := &Weapon{config, config.MainHandWeapon}
	round := CombatRoundSimulationResult{
		attacks: make([]*AttackSimulationResult, attacker.AttacksPerRound()+attacker.ExtraAttacks()),
	}

	addRound := func(roundPos int, attackResult *AttackSimulationResult) {
		round.attacks[roundPos] = attackResult
		round.damage += attackResult.finalDamage
		if attackResult.hit {
			round.hits++
		}
		if attackResult.criticalConfirmed {
			round.crits++
		}
	}

	attackPenaltySchedule := attacker.AttackPenaltySchedule()
	logger.Trace.Println("Attack penalty schedule", attackPenaltySchedule)
	for i := 0; i < attacker.AttacksPerRound()+attacker.ExtraAttacks(); i++ {
		logger.Trace.Println("Simulating round", i+1)
		attackResult := processAttack(i+1, attackPenaltySchedule[i], config, attacker, mhWeapon)
		addRound(i, attackResult)
	}

	return round
}
