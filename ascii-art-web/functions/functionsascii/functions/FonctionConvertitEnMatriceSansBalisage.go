package artweb

import (
	"os"
)

func ConvertitEnMatriceSansBalisage(lettre byte, NomFichier string) []string {
	fichier, err := os.ReadFile(NomFichier)

	if err != nil {
		return nil
	}

	var Matrice []string

	i := 1
	var CompteurAscii byte
	CompteurAscii = 32
	NbLignes := 0
	LenLignes := 0
	ModeRecuperation := false
	var ligne []byte
	for i < len(fichier) {
		if (fichier[i] >= ' ' && fichier[i] <= '~') || fichier[i] == '\n' {
			if NomFichier == "thinkertoy.txt" {
				if i > 1 {
					if fichier[i] == '\n' && fichier[i-2] == '\n' && ModeRecuperation == false {
						CompteurAscii++
					}
				}
			} else {
				if fichier[i] == '\n' && fichier[i-1] == '\n' && ModeRecuperation == false {
					CompteurAscii++
				}
			}

			if CompteurAscii == lettre && ModeRecuperation == false {
				ModeRecuperation = true
				if lettre == ' ' {
					i--
				}

			} else if fichier[i] == '\n' && fichier[i-1] == '\n' && ModeRecuperation == true {
				ModeRecuperation = false
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
