package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Module struct {
    typeName string
    moduleName string
    targetModules []string
    inputStates map[string]bool
    outputState bool
}

type Pulse struct {
    isHigh bool
    target string
}

type StateData struct {
    prevState string
    stepNumber int
    countHigh int
    countLow int
    totalHigh int
    totalLow int
}

func splitCommaSepNums(num_list_str string) []int {
	num_list_trimmed := strings.Trim(num_list_str, " ")
	num_str_list := strings.Split(num_list_trimmed, ",")
	output := []int{}

	for _, num_str := range num_str_list {
		num_int, err := strconv.ParseInt(strings.Trim(num_str, " "), 10, 64)
		if err == nil {
			output = append(output, int(num_int))
		} else {
			fmt.Fprintf(os.Stderr, "Error parsing string as number: '%s' (%s)\n", num_str, num_list_str)
		}
	}

	return output
}

func handleButtonPush(moduleMap map[string]Module) (map[string]Module, int, int) {
    countHigh := 0
    countLow := 0

    pulseQueue := []Pulse{{false, "broadcaster"}}

    for len(pulseQueue) > 0 {
        currentPulse := pulseQueue[0]
        pulseQueue = pulseQueue[1:]

        if currentPulse.isHigh {
            countHigh += 1
        } else {
            countLow += 1
        }


        currentModule, currentModuleExists := moduleMap[currentPulse.target]
        if !currentModuleExists {
            continue
        }

        if currentModule.typeName == "broadcaster" {
            for _,nextModuleName := range currentModule.targetModules {
                pulseQueue = append(pulseQueue, Pulse{currentPulse.isHigh, nextModuleName})
            }
        }
    }

    return moduleMap, countHigh, countLow
}

func main() {
    var filename = flag.String("f", "../inputs/d20.sample1.txt", "file to use")
    flag.Parse()
    dat, err := os.ReadFile(*filename)
    if err != nil {
        log.Fatalf("unable to read file: %f", err)
    }
    
    total := 0

    moduleMap := map[string]Module{}
    moduleNameList := []string{}

    p_module := regexp.MustCompile("([%&]?)([A-Za-z]+) -> ([^\t\r\n]*)")

    // Parse file.
    for _,fLine := range strings.Split(strings.Trim(string(dat), " \n"), "\n") {
        moduleSplits := p_module.FindStringSubmatch(fLine)

        destinations := strings.Split(moduleSplits[3], ", ")

        var currentModule Module
        if moduleSplits[1] == "broadcaster" {
            currentModule = Module{"broadcaster", "broadcaster", destinations, map[string]bool{}, false}
        } else {
            currentModule = Module{moduleSplits[1], moduleSplits[2], destinations, map[string]bool{}, false}
        }

        moduleMap[currentModule.moduleName] = currentModule
        moduleNameList = append(moduleNameList, currentModule.moduleName)
    }

    // Setup Conjunction links.
    for _,currentModName := range moduleNameList {
        for _,destModName := range moduleMap[currentModName].targetModules {
            destMod := moduleMap[destModName]

            if destMod.typeName == "&" {
                moduleMap[destModName].inputStates[currentModName] = false
            }
        }
    }


    //stateMap := map[string]StateData{}


    //fmt.Printf("node_list: %v\n", node_list)
    fmt.Printf("T: %d\n", total)
}

