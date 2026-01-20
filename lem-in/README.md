# lem-in


## 🐜 Overview
`lem-in` is a path-finding project that simulates ants moving through a network of rooms.  
The goal is to send all ants from the **start** room to the **end** room in the **fewest possible turns**, following strict movement rules.


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

- input file path is read from `stdin`
- output : step by step movement of ants
- on invalid input file, prints `ERROR: invalid data format`, followed by more precise error source


## 👤 Authors 

Oscar Hernandez-Jaussely - Project Manager & Input Treatment
<br>
Hocine Belagha - Paths Handling
<br>
Romain Maillard - Flow Distribution
