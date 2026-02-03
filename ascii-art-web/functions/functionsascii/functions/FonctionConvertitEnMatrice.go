package artweb

import (
	"os"
)

func ConvertitEnMatrice(lettre byte, NomFichier string) []string {
	fichier, err := os.ReadFile(NomFichier)

	if err != nil {
		return nil
	}

	var Matrice []string

	i := 0
	NbLignes := 0
	LenLignes := 0
	ModeRecuperation := false
	var ligne []byte
	for i < len(fichier) {
		if (fichier[i] >= ' ' && fichier[i] <= '~') || fichier[i] == '\n' {
			if fichier[i] == lettre && fichier[i+1] == '\n' && ModeRecuperation == false {
				ModeRecuperation = true
				i++
			} else if fichier[i] == '.' && fichier[i+1] == lettre && ModeRecuperation == true {
				ModeRecuperation = false
				i++
				break
			} else if ModeRecuperation == true {
				if fichier[i] == '\n' {
					if ligne != nil {
						if len(Matrice) == 0 {
							LenLignes = len(ligne)
						} else {
							if len(ligne) != LenLignes {
								break
							}
						}
						Matrice = append(Matrice, string(ligne))
						NbLignes++
						ligne = nil
					} else {
						break
					}
				} else {
					ligne = append(ligne, fichier[i])
				}
			}
		}
		i++
	}

	if len(Matrice) == 8 {
		return Matrice
	} else {
		return nil
	}
}
