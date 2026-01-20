package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	AntNumber int
	StartRoom string
	EndRoom   string
)

type Room struct {
	id        string
	neighbors []*Room
	visited   bool
	coordX    int
	coordY    int
}

type Labyrinth struct {
	Rooms map[string]*Room
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("ERROR : Missing input file")
		return
	}

	filename := os.Args[1]
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("ERROR : Impossible to read input file")
		return
	}
	fileContent := strings.Trim(string(file), "\n")
	fileLines := strings.Split(fileContent, "\n")

	// set ant number if present and not 0
	tmp, err := strconv.Atoi(fileLines[0])
	if err != nil || tmp == 0 {
		fmt.Println("ERROR: invalid data format, invalid number of Ants")
		return
	}
	AntNumber = tmp

	// check for presence of start & end rooms
	// + check for correct position in input file
	start, end := false, false
	for x, line := range fileLines {
		if line == "##start" && (x != len(fileLines)-1) && (len(strings.Split(fileLines[x+1], " ")) == 3) {
			start = true
		}
		if line == "##end" && (x != len(fileLines)-1) && (len(strings.Split(fileLines[x+1], " ")) == 3) {
			end = true
		}
	}
	if !start {
		fmt.Println("ERROR: invalid data format, no start room found")
		return
	}
	if !end {
		fmt.Println("ERROR: invalid data format, no end room found")
		return
	}

	// create rooms instance if valid //

	LabyrinthRooms := make(map[string]*Room)
	Labyrinth := Labyrinth{
		Rooms: LabyrinthRooms,
	}

	startNext, endNext := false, false
	for x, line := range fileLines {
		// skip ant number
		if x == 0 {
			continue
		}
		if lineSlice := strings.Split(line, " "); len(lineSlice) == 3 {
			// ROOM
			id := lineSlice[0]

			if id[0] == 'L' || id[0] == 'l' {
				fmt.Printf("ERROR: invalid data format, room \"%v\" has invalid name\n", id)
				return
			}

			coordX, err1 := strconv.Atoi(lineSlice[1])
			coordY, err2 := strconv.Atoi(lineSlice[2])
			if err1 != nil || err2 != nil {
				fmt.Printf("ERROR: invalid data format, room %v has invalid coordinates\n", id)
				return
			}

			if startNext {
				StartRoom = id
				startNext = false
			}
			if endNext {
				EndRoom = id
				endNext = false
			}

			room := Room{
				id:      id,
				visited: false,
				coordX:  coordX,
				coordY:  coordY,
			}
			LabyrinthRooms[id] = &room
		} else if lineSlice := strings.Split(line, "-"); len(lineSlice) == 2 {
			// PATH
			if lineSlice[0] == lineSlice[1] {
				fmt.Printf("ERROR: invalid data format, room %v is neighboring itself\n", lineSlice[0])
				return
			}
			for _, id := range lineSlice {
				neighbors := []*Room{}
				for _, otherId := range lineSlice {
					// skip own id
					if otherId == id {
						continue
					}
					neighbors = append(neighbors, Labyrinth.Rooms[otherId])
				}
				Labyrinth.Rooms[id].neighbors = append(Labyrinth.Rooms[id].neighbors, neighbors...)
			}
		} else if line[0] == '#' {
			// comment or start/end

			if line == "##start" {
				startNext = true
			}
			if line == "##end" {
				endNext = true
			}
		} else {
			fmt.Println("ERROR: invalid data format")
			return
		}
	}

	allPaths := FindAllPaths(&Labyrinth)
	if allPaths == nil {
		fmt.Println("ERROR: invalid data format, no path found from start to end")
		return
	}

	rawPaths := FilterPathsNoSharedRooms(allPaths, true)
	paths := ToPaths(rawPaths)
	result := DistributeAnts(paths, AntNumber)
	resultByDeparture := SortByDeparture(result)
	movements := GenerateMovements(resultByDeparture)

	// print output
	for _, line := range fileLines {
		fmt.Println(line)
	}
	fmt.Println()
	for _, turn := range movements {
		fmt.Println(strings.Join(turn, " "))
	}
}

func getTurns(result [][]string) int {
	x := 0
	for _, ant := range result {
		if len(ant) > x {
			x = len(ant)
		}
	}
	return x
}
