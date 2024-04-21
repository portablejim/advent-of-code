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
    source string
}

type StateData struct {
    nextState string
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

    pulseQueue := []Pulse{{false, "broadcaster", "button"}}

    for len(pulseQueue) > 0 {
        currentPulse := pulseQueue[0]
        pulseQueue = pulseQueue[1:]

        //pulseLevelStr := "L"
        if currentPulse.isHigh {
            countHigh += 1
            //pulseLevelStr = "H"
        } else {
            countLow += 1
        }


        currentModule, currentModuleExists := moduleMap[currentPulse.target]
        if !currentModuleExists {
            //fmt.Printf("P I: (from %s) %s -%s-> ?!?!\n", currentPulse.source, currentPulse.target, pulseLevelStr)
            continue
        }
        //fmt.Printf("P I: %s -%s-> %s %v\n", currentPulse.source, pulseLevelStr, currentPulse.target, currentModule.targetModules)

        if currentModule.moduleName == "button" || currentModule.moduleName == "broadcaster" {
            for _,nextModuleName := range currentModule.targetModules {
                //fmt.Printf("P O: %s -%s-> %s\n", currentPulse.target, pulseLevelStr, nextModuleName)
                pulseQueue = append(pulseQueue, Pulse{currentPulse.isHigh, nextModuleName, currentModule.moduleName})
            }
        } else if currentModule.typeName == "%" {
            // Flip-flop
            if !currentPulse.isHigh {
                currentModule.outputState = !currentModule.outputState
                moduleMap[currentPulse.target] = currentModule
                for _,nextModuleName := range currentModule.targetModules {
                    /*
                    outputPulseLevelStr := "L"
                    if currentModule.outputState {
                        outputPulseLevelStr = "H"
                    }
                    fmt.Printf("P O: %s -%s-> %s\n", currentPulse.target, outputPulseLevelStr, nextModuleName)
                    */
                    pulseQueue = append(pulseQueue, Pulse{currentModule.outputState, nextModuleName, currentModule.moduleName})
                }
            } else {
                //fmt.Printf("P O: %s END (%d)\n", currentPulse.target, len(pulseQueue))
            }
        } else if currentModule.typeName == "&" {
            // Conjunction
            currentModule.inputStates[currentPulse.source] = currentPulse.isHigh

            areAllInputsHigh := true
            for _,inputIsHigh := range currentModule.inputStates {
                areAllInputsHigh = areAllInputsHigh && inputIsHigh
            }

            for _,nextModuleName := range currentModule.targetModules {
                /*
                outputPulseLevelStr := "L"
                if !areAllInputsHigh {
                    outputPulseLevelStr = "H"
                }
                fmt.Printf("P O: %s -%s-> %s\n", currentPulse.target, outputPulseLevelStr, nextModuleName)
                */
                pulseQueue = append(pulseQueue, Pulse{!areAllInputsHigh, nextModuleName, currentModule.moduleName})
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

    stateMap := map[string]StateData{}

    totalHigh := 0
    totalLow := 0

    currentState := fmt.Sprintf("%v", moduleMap)
    stateMap[currentState] = StateData{"", 0, 0, 0, 0, 0}

    loopSize := -1

    pressCount := 1_000
    for i := 0; i < pressCount; i += 1 {
        //fmt.Printf("I: %d %v\n", i, currentState)
        currentStateData, currentStateInMap := stateMap[currentState]
        if currentStateInMap && currentStateData.nextState != "" {
            if loopSize < 0 {
                loopSize = i - currentStateData.stepNumber
                fmt.Printf("Loop: %d %d\n", i, currentStateData.stepNumber)

                numRemaining := pressCount - i
                loopCount := (numRemaining / loopSize)
                numRemainAfter := numRemaining % loopSize
                loopTotalHigh := totalHigh - currentStateData.totalHigh - currentStateData.countHigh
                loopTotalLow := totalLow - currentStateData.totalLow

                totalHigh += loopCount * loopTotalHigh
                totalLow += loopCount * loopTotalLow
                i = pressCount - numRemainAfter
                //fmt.Printf("Loop found %d %d %d\n", loopSize, numRemaining, numRemainAfter)
            } else {
                //fmt.Printf("Loop used\n")
            }
            currentState = currentStateData.nextState


        } else {
            _, pushHighCount, pushLowCount := handleButtonPush(moduleMap)

            totalHigh += pushHighCount
            totalLow += pushLowCount

            //fmt.Printf("State: %s\n", currentState)
            nextState := fmt.Sprintf("%v", moduleMap)


            if currentStateData.nextState == "" {
                currentStateData.nextState = nextState
                stateMap[currentState] = currentStateData;
            }

            nextStateData, nextStateInMap := stateMap[nextState]
            if !nextStateInMap {
                fmt.Printf("Current state not in map\n")
                nextStateData = StateData{"", i+1, pushHighCount, pushLowCount, totalHigh, totalLow}
                stateMap[nextState] = nextStateData
            }
            currentState = nextState

            //fmt.Printf("State: %s\n", currentState)
        }
    }

    total = totalLow * totalHigh

    //fmt.Printf("node_list: %v\n", node_list)
    fmt.Printf("T: %d %d %d\n", totalHigh, totalLow, total)
}

