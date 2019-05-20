package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"simarelith/logger"
	"simarelith/simulation"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/ini.v1"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	configPath := flag.String("c", "", "The path to the config file")
	loggingLevel := flag.Int("l", 1, "Logging level: 0:None, 1:Error, 2:Warning, 3:Info, 4:Trace")
	varyAc := flag.Int("ac", 0, "Re-run the simulation this many times adding 1 to the target ac every time")
	iterations := flag.Int("i", 100, "How many iterations of the sim to run")
	outputFile := flag.String("o", "", "The path to the output file")

	flag.Parse()

	logger.Init(*loggingLevel)

	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	logger.Trace.Println("Reading config at", *configPath)
	configAbsPath, _ := filepath.Abs(*configPath)
	cfg, err := ini.Load(configAbsPath)
	if err != nil {
		fmt.Printf("Fail to read config file: %v\n", err)
		os.Exit(1)
	}

	logger.Trace.Println("Read successful, parsing config")
	simConfig, _ := simulation.ConfigFromIni(cfg)
	logger.Trace.Println("Parse successful, running simulation")
	simulator := simulation.NewSimulator(*simConfig)

	rounds := simulator.Run(*iterations)
	analysis := simulation.NewAnalysis(rounds)
	logger.Info.Println(spew.Sdump(analysis))

	outputFilePath, _ := filepath.Abs(*outputFile)
	output, err := os.Create(outputFilePath)
	defer output.Close()
	if err != nil {
		logger.Error.Println("Could not create file", outputFilePath)
		os.Exit(1)
	}

	writer := csv.NewWriter(output)
	defer writer.Flush()

	writer.Write([]string{
		"AC",
		"TotalAttacks",
		"FirstHitPercentage",
		"HitPercentage",
		"CritPercentage",
		"TotalDamage",
		"DamagePerRound",
	})

	for i := 0; i <= *varyAc; i++ {
		logger.Info.Println("Simming at target AC", simConfig.Target.ArmorClass)
		rounds := simulator.Run(*iterations)
		analysis := simulation.NewAnalysis(rounds)

		writer.Write(
			[]string{
				strconv.Itoa(simConfig.Target.ArmorClass),
				strconv.Itoa(analysis.TotalAttacks),
				fmt.Sprintf("%f", analysis.FirstHitPercentage),
				fmt.Sprintf("%f", analysis.HitPercentage),
				fmt.Sprintf("%f", analysis.CritPercentage),
				strconv.Itoa(analysis.TotalDamage),
				fmt.Sprintf("%f", analysis.DamagePerRound),
			},
		)

		simConfig.Target.ArmorClass++
	}
}
