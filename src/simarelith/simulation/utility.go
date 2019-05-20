package simulation

import (
	"fmt"
	"strconv"

	"github.com/justinian/dice"
)

func roll(description string) int {
	if res, err := strconv.Atoi(description); err == nil {
		return res
	}

	res, _, err := dice.Roll(description)
	if err != nil {
		fmt.Println("Error rolling dice with description: ", description)
		panic("Stopping simulation")
	}
	return res.Int()
}
