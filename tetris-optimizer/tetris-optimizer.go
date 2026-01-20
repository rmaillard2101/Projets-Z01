package main

import (
	"fmt"
	"os"
)

func main() {
	fichier, err := os.ReadFile("sample.txt")
	if err != nil {
		fmt.Println("ERROR")
		os.Exit(1)
	}
	var tetrotab [][][]byte
	var tetro [][]byte
	var tetroligne []byte

	//capture tetros
	nbcarres := 0
	for i := 0; i < len(fichier); i++ {
		if fichier[i] != '.' && fichier[i] != '#' && fichier[i] != '\n' {
			fmt.Println("ERROR")
			os.Exit(1)
		}
		if fichier[i] == '.' || fichier[i] == '#' {
			if fichier[i] == '#' {
				nbcarres++
			}
			tetroligne = append(tetroligne, fichier[i])
		} else if fichier[i] == '\n' && len(tetroligne) > 0 {
			tetro = append(tetro, tetroligne)
			tetroligne = nil
		} else if fichier[i] == '\n' && len(tetroligne) == 0 {
			if nbcarres != 4 {
				fmt.Println("ERROR")
				os.Exit(1)
			}
			nbcarres = 0
			if tetro != nil {
				if len(tetro) != 4 || len(tetro[0]) != 4 {
					fmt.Println("ERROR")
					os.Exit(1)
				}
				if !IsValidTetrimino(tetro) {
					fmt.Println("ERROR")
					os.Exit(1)
				}
				tetrotab = append(tetrotab, tetro)
				tetro = nil
			}
		}
		if fichier[i] == '\n' && i == len(fichier)-1 {
			if nbcarres != 4 {
				fmt.Println("ERROR")
				os.Exit(1)
			}
			nbcarres = 0
			if tetro != nil {
				if len(tetro) != 4 || len(tetro[0]) != 4 {
					fmt.Println("ERROR")
					os.Exit(1)

				}
				if !IsValidTetrimino(tetro) {
					fmt.Println("ERROR")
					os.Exit(1)
				}
				tetrotab = append(tetrotab, tetro)
				tetro = nil
			}
		}
	}

	if len(tetrotab) == 0 {
		fmt.Println("ERROR")
		os.Exit(1)
	}
	tetrovars := make([][]int, len(tetrotab))
	for i := 0; i < len(tetrovars); i++ {
		tetrovars[i] = make([]int, 2)
	}

	//justification tetros
	for i := 0; i < len(tetrotab); i++ {
		limxi := 0
		limximodif := false
		limxe := 0
		limxemodif := false
		limyi := 0
		limyimodif := false
		limye := 0
		limyemodif := false
		for j := 0; j < len(tetrotab[i]); j++ {
			LigneVide := true
			SurDiese := false
			for k := 0; k < len(tetrotab[i][j]); k++ {
				if tetrotab[i][j][k] == '#' {
					LigneVide = false
					if SurDiese == false {
						if limximodif == false {
							limxi = k
							limximodif = true
						} else if limximodif == true {
							if limxi > k {
								limxi = k
							}
						}

					}
					SurDiese = true
					if k == len(tetrotab[i][j])-1 && SurDiese == true {
						limxe = k
						limxemodif = true
					}
				} else if tetrotab[i][j][k] == '.' {
					if SurDiese == true {
						if limxemodif == false {
							limxe = k - 1
							limxemodif = true
						} else if limxemodif == true {
							if limxe < k-1 {
								limxe = k - 1
							}
						}
					}
					SurDiese = false
				}
			}
			if LigneVide == false && limyimodif == false {
				limyi = j
				limyimodif = true
			} else if LigneVide == true && limyimodif == true && limyemodif == false {
				limye = j - 1
				limyemodif = true
			}
			if len(tetrotab[i]) > 0 {
				if j == len(tetrotab[i])-1 && limyemodif == false {
					limye = j
					limyemodif = true
				}
			}

		}

		//production damier

		var temptetro [][]byte
		var temptetroligne []byte
		for j := limyi; j <= limye; j++ {
			for k := limxi; k <= limxe; k++ {
				temptetroligne = append(temptetroligne, tetrotab[i][j][k])
			}
			temptetro = append(temptetro, temptetroligne)
			temptetroligne = nil
		}

		tetrotab[i] = nil
		for j := 0; j < len(temptetro); j++ {
			for k := 0; k < len(temptetro[j]); k++ {
				temptetroligne = append(temptetroligne, temptetro[j][k])
			}
			tetrotab[i] = append(tetrotab[i], temptetroligne)
			temptetroligne = nil
		}
	}

	var damier [][]byte
	var lignevide []byte
	largeur := 1
	nbtetro := 0
	for i := 0; i < largeur; i++ {
		lignevide = append(lignevide, '.')
	}

	for i := 0; i < largeur; i++ {
		ligne := make([]byte, len(lignevide))
		copy(ligne, lignevide)
		damier = append(damier, ligne)
	}

	for nbtetro < len(tetrotab) {
		posable, x, y := Placeable(&damier, &tetrotab[nbtetro], tetrovars[nbtetro][0], tetrovars[nbtetro][1])
		if posable == true {
			Ecriture(&damier, &tetrotab[nbtetro], x, y, 'A'+byte(nbtetro))
			tetrovars[nbtetro][0] = x
			tetrovars[nbtetro][1] = y
			nbtetro++
		} else if posable == false {
			if nbtetro > 0 {
				tetrovars[nbtetro][0] = 0
				tetrovars[nbtetro][1] = 0
				nbtetro--
				Ecriture(&damier, &tetrotab[nbtetro], tetrovars[nbtetro][0], tetrovars[nbtetro][1], '.')
				if tetrovars[nbtetro][0] < largeur-1 {
					tetrovars[nbtetro][0]++
				} else if tetrovars[nbtetro][0] == largeur-1 {
					tetrovars[nbtetro][0] = 0
					tetrovars[nbtetro][1]++
				}

			} else if nbtetro == 0 {
				largeur++
				damier = nil
				lignevide = nil
				for i := 0; i < largeur; i++ {
					lignevide = append(lignevide, '.')
				}
				for i := 0; i < largeur; i++ {
					ligne := make([]byte, len(lignevide))
					copy(ligne, lignevide)
					damier = append(damier, ligne)
				}
			}
		}
	}

	//impression
	for i := 0; i < len(tetrotab); i++ {
		for j := 0; j < len(tetrotab[i]); j++ {
			for k := 0; k < len(tetrotab[i][j]); k++ {
				fmt.Printf(string(tetrotab[i][j][k]))
				if k == len(tetrotab[i][j])-1 {
					fmt.Printf("\n")
				}
			}
			if j == len(tetrotab[i])-1 {
				fmt.Printf("\n")
			}
		}
	}

	for i := 0; i < len(damier); i++ {
		var printligne []byte
		for j := 0; j < len(damier[i]); j++ {
			printligne = append(printligne, damier[i][j])
		}
		fmt.Println(string(printligne))
	}
}

func Placeable(damieraddr *[][]byte, tetro *[][]byte, bordx int, bordy int) (bool, int, int) {
	var valide bool
	var x int
	var y int
	for i := bordy; i < len(*damieraddr)-(len(*tetro)-1); i++ {
		for j := bordx; j < len((*damieraddr)[i])-(len((*tetro)[0])-1); j++ {
			superposition := false
			for k := 0; k < len(*tetro); k++ {
				for l := 0; l < len((*tetro)[k]); l++ {
					if (*tetro)[k][l] == '#' {
						if (*damieraddr)[i+k][j+l] >= 'A' && (*damieraddr)[i+k][j+l] <= 'Z' {
							superposition = true
						}
					}
				}
			}
			if superposition == false {
				valide = true
				x = j
				y = i
				goto retour
			}
		}
		bordx = 0
	}
retour:
	return valide, x, y
}

func Ecriture(damieraddr *[][]byte, tetro *[][]byte, bordx int, bordy int, lettre byte) {
	for k := 0; k < len(*tetro); k++ {
		for l := 0; l < len((*tetro)[k]); l++ {
			if (*tetro)[k][l] == '#' {
				(*damieraddr)[bordy+k][bordx+l] = lettre
			}
		}
	}
}

func IsValidTetrimino(tetro [][]byte) bool {
	visited := make([][]bool, len(tetro))
	for i := range visited {
		visited[i] = make([]bool, len(tetro[i]))
	}

	var startY, startX int
	found := false

	// Trouve un premier '#' pour démarrer la recherche
	for y := 0; y < len(tetro); y++ {
		for x := 0; x < len(tetro[y]); x++ {
			if tetro[y][x] == '#' {
				startY = y
				startX = x
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return false
	}

	// Lance la recherche pour marquer les cases connectées
	dfsMark(tetro, visited, startY, startX)

	// Vérifie s'il reste des '#' non visités (donc non connectés)
	for y := 0; y < len(tetro); y++ {
		for x := 0; x < len(tetro[y]); x++ {
			if tetro[y][x] == '#' && !visited[y][x] {
				return false // Ce '#' est isolé
			}
		}
	}
	return true
}

func dfsMark(tetro [][]byte, visited [][]bool, y, x int) {
	if y < 0 || y >= len(tetro) || x < 0 || x >= len(tetro[y]) {
		return
	}
	if tetro[y][x] != '#' || visited[y][x] {
		return
	}

	visited[y][x] = true

	dfsMark(tetro, visited, y-1, x)
	dfsMark(tetro, visited, y+1, x)
	dfsMark(tetro, visited, y, x-1)
	dfsMark(tetro, visited, y, x+1)
}
