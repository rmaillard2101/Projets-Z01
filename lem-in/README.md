[GitHub deposit link (Oscar Hernandez)](https://github.com/IronBeagle404/lem-in)


# lem-in

## 🐜 Overview

`lem-in` is a path-finding project that simulates ants moving through a network of rooms.  
The goal is to send all ants from the **start** room to the **end** room in the **fewest possible turns**, following strict movement rules.

## 🐜 Ant Movement Rules

-Each ant follows strict movement constraints during the simulation:
-Each ant can move only once per turn
-An ant can only move through one tunnel per turn
-Each tunnel can be used by only one ant per turn
-Each room can contain only one ant at a time
-Except ##start and ##end, which can contain any number of ants
-Ants cannot move into an occupied room
-Ants cannot cross or swap positions
-Only ants that actually move during a turn are displayed in the output

These rules ensure a realistic simulation and enforce the need for optimal path selection and traffic management.

## ⚙️ Features

- Parses an input describing:
  - Number of ants
  - Rooms (with coordinates)
  - Links between rooms
- Detects invalid input and prints detailed error
- Computes optimal paths (using BFS / max-flow logic)
- Simulates and prints ant movements in the format :
  `L1-roomA L2-roomB`

## 🚀 Usage

```
go run . input.txt
```

The program reads the input file from stdin and outputs the full simulation to stdout.

- 📥 Input Format

  The input describes the ant colony and follows this structure:
  -Number of ants (positive integer)
  -Room definitions
  -Tunnel links

  -Rooms
  Defined as: room_name x y
  room_name:
  Cannot start with L or #
  Contains no spaces
  Coordinates x and y are integers

  -Special commands:
  ##start → next room is the start room
  ##end → next room is the end room
  Lines starting with # are treated as comments and ignored

  -Links
  Defined as: room1-room2
  A tunnel connects exactly two rooms
  No duplicate tunnels
  No self-links (room1-room1)
  Links to unknown rooms are considered invalid

- 📤 Output Format

  The output is composed of two parts:
  -1. Echo of the input
  The program first prints the validated input data:
  -number of ants
  -rooms
  -links
  -2. Ant movements
  -Each line represents one turn

  Format: L<ant_number>-<room_name>

  Ant numbers go from 1 to the number of ants

  Example:

  ```
  L1-roomA L2-roomB
  L1-roomC L2-roomD L3-roomB
  ```

🗂 Project Structure

```
.
├── main.go # Program entry point
├── parser/ # Input parsing and validation
├── graph/ # Graph representation and traversal
├── solver/ # Path-finding and optimization logic
├── simulator/ # Ant movement simulation
├── errors/ # Custom error handling
└── tests/ # Unit and integration tests
```

## 👤 Authors

Author details are available in the `creditentials.md` file.
