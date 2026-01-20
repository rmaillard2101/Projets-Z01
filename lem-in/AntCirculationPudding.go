package main

import (
	"fmt"
	"sort"
	"strconv"
)

type Path struct {
	Rooms []string
	Len   int
	Stage int
}

func ToPaths(raw [][]string) []Path {
	paths := make([]Path, len(raw))
	for i, s := range raw {
		paths[i] = Path{
			Rooms: s,
			Len:   len(s),
			Stage: 0,
		}
	}
	return paths
}

func DistributeAnts(paths []Path, numAnts int) [][]string {
	antPaths := make([][]string, len(paths))

	for i := 1; i <= numAnts; i++ {
		var shortestPathIdx int
		var shortestLen int

		for j := 0; j < len(paths); j++ {
			if shortestLen == 0 || paths[j].Len+paths[j].Stage < shortestLen {
				shortestPathIdx = j
				shortestLen = paths[j].Len + paths[j].Stage
			}
		}
		paths[shortestPathIdx].Stage++
		antPaths[shortestPathIdx] = append(antPaths[shortestPathIdx], strconv.Itoa(i))
	}

	result := make([][]string, numAnts)

	for antNum := 1; antNum <= numAnts; antNum++ {
		found := false

		for pathIdx, members := range antPaths {
			for pos, f := range members {
				fnum, _ := strconv.Atoi(f)
				if fnum == antNum {
					antSteps := make([]string, 0, pos+len(paths[pathIdx].Rooms))
					for k := 0; k < pos; k++ {
						antSteps = append(antSteps, "#")
					}
					antSteps = append(antSteps, paths[pathIdx].Rooms...)
					result[antNum-1] = antSteps
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}

	return result
}

func SortByDeparture(ants [][]string) [][]string {
	sort.SliceStable(ants, func(i, j int) bool {
		countHash := func(s []string) int {
			c := 0
			for _, v := range s {
				if v == "#" {
					c++
				} else {
					break
				}
			}
			return c
		}
		return countHash(ants[i]) < countHash(ants[j])
	})
	return ants
}

func GenerateMovements(ants [][]string) [][]string {
	maxLen := 0
	for _, steps := range ants {
		if len(steps) > maxLen {
			maxLen = len(steps)
		}
	}

	result := [][]string{}

	for turn := 0; turn < maxLen; turn++ {
		moves := []string{}
		for i, steps := range ants {
			if turn < len(steps) && steps[turn] != "#" {
				moves = append(moves, fmt.Sprintf("L%d-%s", i+1, steps[turn]))
			}
		}
		if len(moves) > 0 {
			result = append(result, moves)
		}
	}

	return result
}
