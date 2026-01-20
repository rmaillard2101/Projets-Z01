package main

import "sort"

func FindAllPaths(labyrinth *Labyrinth) [][]string {
	// Verification de l'existance du labyrinth et des start/endRoom (évite les panic si la structure est vide ou mal initialisée)
	if labyrinth == nil || len(labyrinth.Rooms) == 0 || StartRoom == "" || EndRoom == "" {
		return nil
	}

	// Recupere la start et end room a partir du nom
	startRoom, startExists := labyrinth.Rooms[StartRoom]
	endRoom, endExists := labyrinth.Rooms[EndRoom]

	if !startExists || !endExists {
		return nil
	}

	// Là ou l'on va stocker toutes les routes
	var allPaths [][]string

	// Map pour marque les room déjà visitées dans la route en cours (la fameuse visited que disait romain)
	visitedRooms := make(map[string]bool)

	// On crée une variable qui va parcourir toutes les room voisines et si elle est visited alors on skip pour pas la revisiter
	var explore func(currentRoom *Room, currentPath []string)

	explore = func(currentRoom *Room, currentPath []string) {
		if visitedRooms[currentRoom.id] {
			return
		}

		// Marque la room actuelle comme visited
		visitedRooms[currentRoom.id] = true

		// Ajoute le nom de la room au chemin en cours
		currentPath = append(currentPath, currentRoom.id)

		// Si on atteint la endRoom, cela signifie que le chemin est completer et on fait une copie du chemin dans une variable finalpath et sera ensuite
		// ajoute a allpaths donc la liste de path que l'on renvoie a la fin
		// (On crée la variable finalpath pour y mettre la copie du chemin afin qu'il ne sois pas modifier plus tard)
		if currentRoom == endRoom {
			finalPath := make([]string, len(currentPath))
			copy(finalPath, currentPath)

			allPaths = append(allPaths, finalPath)

			// On passe endroom en non visited avant de commence le retour en arrière jusqu'a start afin de commence la prochaine route (backtracking)
			visitedRooms[currentRoom.id] = false
			return
		}

		// On explore chacune des room voisine
		for _, neighborRoom := range currentRoom.neighbors {
			// Une sécurite au cas on a un chemin entre deux room mais que l'une des deux n'existe pas
			if neighborRoom == nil {
				continue
			}

			// Si la prochaine room n'a pas été visited, on continue d'avance a la recherche d'un chemin jusqu'a end et dans le cas
			// ou le chemin ne mène a rien ou fini par n'avoir que des room deja visited en option, il fait marche arrière pour tenté un autre chemin (le dps)
			if !visitedRooms[neighborRoom.id] {
				explore(neighborRoom, currentPath)
			}
		}

		// Retire les room du visited petit a petit en marche arrière jusqu'a startRoom pour préparer la prochaine route
		visitedRooms[currentRoom.id] = false
	}

	// On repart de la startRoom pour commence le chemin suivant
	explore(startRoom, []string{})

	// On renvoie la liste complète des chemins trouvés
	return allPaths
}

// La fonction filtre la liste de chemins pour garder uniquement ceux qui ne partagent aucune room intermédiaire (hors start et end)
// Elle essaie plusieurs options de tri (ordre croissant, décroissant, original et choisit celle qui demande le moins de tours total
func FilterPathsNoSharedRooms(paths [][]string, ignoreStartEnd bool) [][]string {

	// --- Fonction interne qui applique le filtrage "greedy" selon l'ordre des chemins donnés ---
	doFilter := func(in [][]string) [][]string {

		// Retient les room déjà utilisées par les chemins retenus
		used := make(map[string]bool)

		// Retient les chemins validées qui n'entrent pas en conflit
		res := make([][]string, 0, len(in))

		// Boucle sur tous les chemins proposés
		for _, p := range in {

			// Pour marquer tout les chemins qui entre en conflit
			conflict := false

			// On vérifie chaque room du chemin
			for i, room := range p {
				// On ignore les room start et end
				if ignoreStartEnd && (i == 0 || i == len(p)-1) {
					continue
				}
				// Si la room est déjà utilisée (conflit), on saute le chemin
				if used[room] {
					conflict = true
					break
				}
			}

			// Si on a un conflit, on passe au chemin suivant
			if conflict {
				continue
			}

			// Si le chemin est valide, on l'ajoute au résultat
			res = append(res, p)

			// On marque toutes ses room intermédiaires comme used
			for i, room := range p {
				if ignoreStartEnd && (i == 0 || i == len(p)-1) {
					continue
				}
				used[room] = true
			}
		}

		// Retourne les chemins filtrés (sans conflits)
		return res
	}

	// L'option A qui va trier les chemin du plus court au plus long
	// On crée une copie des chemin pour ensuite les trier et appliquer le filtre greedy dessus
	asc := make([][]string, len(paths))
	copy(asc, paths)
	sort.Slice(asc, func(i, j int) bool { return len(asc[i]) < len(asc[j]) })
	candA := doFilter(asc)

	// L'option B qui va trier les chemin du plus long au plus court
	// On crée une copie des chemin pour ensuite les trier et appliquer le filtre greedy dessus
	desc := make([][]string, len(paths))
	copy(desc, paths)
	sort.Slice(desc, func(i, j int) bool { return len(desc[i]) > len(desc[j]) })
	candB := doFilter(desc)

	// L'option C qui va pas trier les chemin mais juste copier les chemins comme ils sont sans tri et appliquer le filtre greedy dessus
	orig := make([][]string, len(paths))
	copy(orig, paths)
	candC := doFilter(orig)

	// On met l'option A par défaut car c'est celle qui prend de base le plus petit chemin en premier
	// On met une limite de tour initial pratiquement infini qui sera remplace par ceux des options pour la comparaison
	best := candA
	bestTurns := 1 << 30

	// C'est la fonction interne qui va calculer la taille des trois chemins
	tryEvaluate := func(candidate [][]string) int {
		// conversion en structure adaptée à la distribution des fourmis
		// Simulation de la distribution des fourmis (fonction de romain) sur les chemins
		// Calcul du nombre de tours total nécessaire pour le chemin
		p := ToPaths(candidate)
		dist := DistributeAnts(p, AntNumber)
		return getTurns(dist)
	}

	// On teste l'option A
	if t := tryEvaluate(candA); t < bestTurns {
		bestTurns = t
		best = candA
	}

	// On teste l'option B
	if t := tryEvaluate(candB); t < bestTurns {
		bestTurns = t
		best = candB
	}

	// On teste l'option C
	if t := tryEvaluate(candC); t < bestTurns {
		bestTurns = t
		best = candC
	}

	// On va récuperer best et on passe tout dans la variable filtred en ignorant le startRoom pour le retirer du retour
	filtred := make([][]string, 0, len(best))
	for _, p := range best {
		if len(p) == 0 {
			continue
		}
		if p[0] == StartRoom {
			if len(p) > 1 {
				filtred = append(filtred, append([]string(nil), p[1:]...))
			}
		}
	}

	// On retourne la meilleure option de chemins entre A, B et C (sans le Start en tête)
	return filtred
}
