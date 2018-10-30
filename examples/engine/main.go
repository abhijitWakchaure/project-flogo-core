package main

import (
	"flag"
	"fmt"
	"github.com/project-flogo/core/support/log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/project-flogo/core/engine"
)

var cpuProfile = flag.String("cpuprofile", "", "Writes CPU profile to the specified file")
var memProfile = flag.String("memprofile", "", "Writes memory profile to the specified file")

var (
	configProvider engine.AppConfigProvider
)

func main() {

	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create CPU profiling file due to error - %s", err.Error()))
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	e, err := engine.NewFromConfigProvider(configProvider)
	if err != nil {
		log.RootLogger().Errorf("Failed to create engine: %s", err.Error())
		os.Exit(1)
	}

	code := engine.RunEngine(e)

	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create memory profiling file due to error - %s", err.Error()))
			os.Exit(1)
		}

		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Println(fmt.Sprintf("Failed to write memory profiling data to file due to error - %s", err.Error()))
			os.Exit(1)
		}
		f.Close()
	}

	os.Exit(code)
}