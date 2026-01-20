// Dimensions de la grille
const cols = 10;
const rows = 20;

// Définition des pièces Tetris avec id, centre de rotation et pattern 4x4
const tetrisPieces = {
  I: {
    id: 1,
    center: { x: 2, y: 1 },
    pattern: [
      [0, 0, 0, 0],
      [1, 1, 1, 1],
      [0, 0, 0, 0],
      [0, 0, 0, 0],
    ],
  },

  J: {
    id: 2,
    center: { x: 1, y: 1 },
    pattern: [
      [1, 0, 0, 0],
      [1, 1, 1, 0],
      [0, 0, 0, 0],
      [0, 0, 0, 0],
    ],
  },

  L: {
    id: 3,
    center: { x: 1, y: 1 },
    pattern: [
      [0, 0, 1, 0],
      [1, 1, 1, 0],
      [0, 0, 0, 0],
      [0, 0, 0, 0],
    ],
  },

  O: {
    id: 4,
    center: { x: 1, y: 1 },
    pattern: [
      [1, 1, 0, 0],
      [1, 1, 0, 0],
      [0, 0, 0, 0],
      [0, 0, 0, 0],
    ],
  },

  S: {
    id: 5,
    center: { x: 1, y: 1 },
    pattern: [
      [0, 1, 1, 0],
      [1, 1, 0, 0],
      [0, 0, 0, 0],
      [0, 0, 0, 0],
    ],
  },

  Z: {
    id: 6,
    center: { x: 1, y: 1 },
    pattern: [
      [0, 0, 0, 0],
      [1, 1, 0, 0],
      [0, 1, 1, 0],
      [0, 0, 0, 0],
    ],
  },

  T: {
    id: 7,
    center: { x: 1, y: 1 },
    pattern: [
      [0, 1, 0, 0],
      [1, 1, 1, 0],
      [0, 0, 0, 0],
      [0, 0, 0, 0],
    ],
  },
};

const grid = [];

// Création de la grille vide

for (let c = 0; c < cols; c++) {
  const col = [];
  for (let r = 0; r < rows; r++) {
    col.push(0);
  }
  grid.push(col);
}

// Récupération des éléments DOM

const container = document.getElementById("grid");
const lvlValue = document.getElementById("lvlValue");
const scoreValue = document.getElementById("scoreValue");
const playButton = document.getElementById("playButton");
const pauseButton = document.getElementById("pauseButton");
const pauseRestart = document.getElementById("pauseRestart");
const pauseResume = document.getElementById("pauseResume");
const gameOverBox = document.getElementById("gameOverBox");
const nextPreview = document.getElementById("nextPreview");
const timerValue = document.getElementById("timerValue");
const fpsValue = document.getElementById("fpsValue");

// Gestion du clavier pour les commandes

document.addEventListener("keydown", (event) => {
  //if (event.repeat) return;
  if (event.key === "ArrowLeft") command = "left";
  if (event.key === "ArrowRight") command = "right";
  if (event.key === "ArrowUp") command = "up";
  if (event.key === "ArrowDown") command = "down";
  if (event.key === " ") command = "start";
  if (event.key === "Escape") command = "pause";
});

// Boutons du jeu

playButton.addEventListener("click", () => {
  command = "start";
  playButton.blur();
});

pauseButton.addEventListener("click", () => {
  command = "pause";
  pauseButton.blur();
});

pauseRestart.addEventListener("click", () => {
  command = "start";
  pauseRestart.blur();
});

pauseResume.addEventListener("click", () => {
  command = "pause";
  pauseResume.blur();
});

//Corrections textuelles

document.getElementById(
  "pauseOverlay"
).firstElementChild.firstChild.textContent = "PAUSE";

document.getElementById("pauseRestart").textContent = "RESTART";
document.getElementById("pauseResume").textContent = "RESUME";

document.getElementById("gameOverBox").textContent = "GAME OVER";

document.getElementById("pauseButton").textContent = "PAUSE";
document.getElementById("playButton").textContent = "PLAY / RESTART";

document.querySelector("#timer").childNodes[0].textContent = "TIMER: ";
document.querySelector("#lvl").childNodes[0].textContent = "LEVEL: ";
document.querySelector("#score").childNodes[0].textContent = "SCORE: ";

// Fonction utilitaire pour délai (Promise)

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// Création des divs pour la grille principale

for (let r = 0; r < rows; r++) {
  for (let c = 0; c < cols; c++) {
    const div = document.createElement("div");
    div.className = "cell";
    div.dataset.id = grid[c][r];
    container.appendChild(div);
  }
}

let nextCells = [];

// Création des divs pour la prévisualisation de la prochaine pièce

for (let r = 0; r < 4; r++) {
  for (let c = 0; c < 4; c++) {
    const div = document.createElement("div");
    div.className = "nextCell";
    div.dataset.id = 0;
    nextPreview.appendChild(div);
    nextCells.push(div);
  }
}

// Mise à jour de la grille visible depuis le tableau interne

function updateGrid() {
  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      const index = r * cols + c;
      container.children[index].dataset.id = grid[c][r];
    }
  }
  fpsCounter++;
}

// Mise à jour de la prévisualisation de la prochaine pièce

function updateNextPiece() {
  const pattern = tetrisPieces[Object.keys(tetrisPieces)[nextId - 1]].pattern;
  let i = 0;

  for (let r = 0; r < 4; r++) {
    for (let c = 0; c < 4; c++) {
      const srcRow = r - 1;
      const cell = nextCells[i];

      let value = 0;
      if (srcRow >= 0 && pattern[srcRow][c]) {
        value = nextId;
      }

      cell.dataset.id = value;
      i++;
    }
  }
}

let command = null;

let fallPromise = null;
let fallDone = false;
let fallDelay = 1000;

// Gestion de la chute automatique

function startFallTimer() {
  fallDone = false;
  fallPromise = delay(fallDelay).then(() => true);

  (async () => {
    while (!fallDone) {
      const done = await Promise.race([fallPromise, delay(5)]);
      if (done === true) fallDone = true;
    }
  })();
}

let playTimerRunning = false;
let playTimerStart = 0;
let playTimerOffset = 0;
let timerPaused = false;

// Gestion du timer de jeu

async function startPlayTimer() {
  playTimerStart = Date.now();
  playTimerOffset = 0;
  playTimerRunning = true;

  while (playTimerRunning) {
    if (!timerPaused) {
      const elapsed = playTimerOffset + (Date.now() - playTimerStart);
      const totalSeconds = Math.floor(elapsed / 1000);
      const hours = Math.floor(totalSeconds / 3600);
      const minutes = Math.floor((totalSeconds % 3600) / 60);
      const seconds = totalSeconds % 60;

      timerValue.textContent =
        String(hours).padStart(2, "0") +
        ":" +
        String(minutes).padStart(2, "0") +
        ":" +
        String(seconds).padStart(2, "0");
    }

    await delay(200);
  }
}

// FPS affichage

let fpsCounter = 0;
let currentFps = 0;

async function displayFps() {
  while (true) {
    const tickPromise = delay(1000);

    currentFps = fpsCounter;
    fpsCounter = 0;
    fpsValue.textContent = currentFps;

    await tickPromise;
  }
}

displayFps();

let spawnX = Math.floor(cols / 2);
let spawnY = 1;
let action = "wait";
let previousAction = null;

let x;
let y;
let id;
let nextId;
let piece;

let lvl;
let score;

// Boucle principale du jeu

function processus() {
  // Commandes start/pause
  if (command === "start") {
    action = "start";
    previousAction = null;
    command = null;
  }

  if (command === "pause") {
    if (action !== "pause") {
      previousAction = action;
      action = "pause";
    } else {
      action = previousAction;
      previousAction = null;
    }
    command = null;
  }

  // Initialisation du jeu au start
  if (action === "start") {
    gameOverBox.style.visibility = "hidden";

    playTimer = 0;
    startPlayTimer();

    // Réinitialisation grille et variables
    for (let c = 0; c < cols; c++) {
      for (let r = 0; r < rows; r++) {
        grid[c][r] = 0;
      }
    }

    x = spawnX;
    y = spawnY;
    id = Math.floor(Math.random() * 7) + 1;
    nextId = Math.floor(Math.random() * 7) + 1;

    lvl = 1;
    score = 0;

    fallPromise = null;
    startFallTimer();

    // Placer la pièce initiale sur la grille
    const keys = Object.keys(tetrisPieces);
    piece = structuredClone(tetrisPieces[keys[id - 1]]);
    const pattern = piece.pattern;
    const cx = piece.center.x;
    const cy = piece.center.y;

    for (let r = 0; r < 4; r++) {
      for (let c = 0; c < 4; c++) {
        if (pattern[r][c] === 1) {
          const gx = x + (c - cx);
          const gy = y + (r - cy);

          if (gx >= 0 && gx < cols && gy >= 0 && gy < rows) {
            grid[gx][gy] = id;
          }
        }
      }
    }
    updateGrid();
    updateNextPiece();
    action = "play";
  }

  //Gestion pause
  if (action === "pause") {
    if (!timerPaused) {
      playTimerOffset += Date.now() - playTimerStart;
      timerPaused = true;
      gameOverBox.style.visibility = "hidden";
      pauseOverlay.style.visibility = "visible";
    }
  } else {
    if (timerPaused) {
      playTimerStart = Date.now();
      timerPaused = false;
      pauseOverlay.style.visibility = "hidden";
    }
  }

  // Gestion du jeu en cours
  if (action === "play") {
    if (command === "left") {
      let canMove = true;

      //Calcul possibilité mouvement
      for (let r = 0; r < 4; r++) {
        for (let c = 0; c < 4; c++) {
          if (piece.pattern[r][c] === 1) {
            const gx = x + (c - piece.center.x);
            const gy = y + (r - piece.center.y);
            const gxL = gx - 1;

            if (gxL < 0 || grid[gxL][gy] !== 0) {
              canMove = false;
            }
            break;
          }
        }
        if (!canMove) break;
      }

      //Effacement réécriture pièce
      if (canMove) {
        for (let r = 0; r < 4; r++)
          for (let c = 0; c < 4; c++)
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              grid[gx][gy] = 0;
            }

        x--;

        for (let r = 0; r < 4; r++)
          for (let c = 0; c < 4; c++)
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              grid[gx][gy] = id;
            }
      }

      command = null;
    }

    if (command === "right") {
      let canMove = true;

      //Calcul possibilité mouvement
      for (let r = 0; r < 4; r++) {
        for (let c = 3; c >= 0; c--) {
          if (piece.pattern[r][c] === 1) {
            const gx = x + (c - piece.center.x);
            const gy = y + (r - piece.center.y);
            const gxR = gx + 1;

            if (gxR >= cols || grid[gxR][gy] !== 0) {
              canMove = false;
            }
            break;
          }
        }
        if (!canMove) break;
      }

      //Effacement réécriture pièce
      if (canMove) {
        for (let r = 0; r < 4; r++)
          for (let c = 0; c < 4; c++)
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              grid[gx][gy] = 0;
            }

        x++;

        for (let r = 0; r < 4; r++)
          for (let c = 0; c < 4; c++)
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              grid[gx][gy] = id;
            }
      }

      command = null;
    }

    if (command === "up") {
      if (id !== 4) {
        const newPattern = [
          [0, 0, 0, 0],
          [0, 0, 0, 0],
          [0, 0, 0, 0],
          [0, 0, 0, 0],
        ];

        let canRotate = true;

        //Effacement piece et calcul possibilité mouvement
        for (let r = 0; r < 4; r++)
          for (let c = 0; c < 4; c++)
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              if (gx >= 0 && gx < cols && gy >= 0 && gy < rows)
                grid[gx][gy] = 0;
            }

        if (id === 1) {
          for (let r = 0; r < 4; r++)
            for (let c = 0; c < 4; c++)
              if (piece.pattern[r][c] === 1) {
                const nr = c;
                const nc = 3 - r;
                newPattern[nr][nc] = 1;
              }
        } else {
          for (let r = 0; r < 4; r++)
            for (let c = 0; c < 4; c++)
              if (piece.pattern[r][c] === 1) {
                const rx = c - piece.center.x;
                const ry = r - piece.center.y;
                const nx = -ry;
                const ny = rx;
                const nr = ny + piece.center.y;
                const nc = nx + piece.center.x;

                if (nr < 0 || nr >= 4 || nc < 0 || nc >= 4) {
                  canRotate = false;
                } else {
                  newPattern[nr][nc] = 1;
                }
              }
        }

        //Possible rotation du pattern puis réécriture pièce
        if (canRotate) {
          for (let r = 0; r < 4; r++) {
            for (let c = 0; c < 4; c++) {
              if (newPattern[r][c] === 1) {
                const gx = x + (c - piece.center.x);
                const gy = y + (r - piece.center.y);
                if (
                  gx < 0 ||
                  gx >= cols ||
                  gy < 0 ||
                  gy >= rows ||
                  grid[gx][gy] !== 0
                ) {
                  canRotate = false;
                  break;
                }
              }
            }
            if (!canRotate) break;
          }
        }

        if (canRotate) {
          piece.pattern = newPattern;
        }

        for (let r = 0; r < 4; r++)
          for (let c = 0; c < 4; c++)
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              if (gx >= 0 && gx < cols && gy >= 0 && gy < rows)
                grid[gx][gy] = id;
            }
      }

      command = null;
    }

    if (command === "down") {
      let canMove = true;

      //Calcul possibilité mouvement
      for (let c = 0; c < 4; c++) {
        for (let r = 3; r >= 0; r--) {
          if (piece.pattern[r][c] === 1) {
            const gx = x + (c - piece.center.x);
            const gy = y + (r - piece.center.y);
            const gyD = gy + 1;

            if (gyD >= rows || grid[gx][gyD] !== 0) {
              canMove = false;
            }
            break;
          }
        }
        if (!canMove) break;
      }

      //Effacement réécriture pièce
      if (canMove) {
        for (let r = 0; r < 4; r++)
          for (let c = 0; c < 4; c++)
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              grid[gx][gy] = 0;
            }

        y++;

        for (let r = 0; r < 4; r++)
          for (let c = 0; c < 4; c++)
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              grid[gx][gy] = id;
            }
      } else {
        fallPromise = Promise.resolve(true);
      }

      command = null;
    }

    //Dans le cas ou le delai
    if (fallPromise) {
      if (fallDone === true) {
        fallDone = false;
        let canFall = true;

        //Processus identique à fleche du bas
        for (let c = 0; c < 4; c++) {
          for (let r = 3; r >= 0; r--) {
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              const gyD = gy + 1;

              if (gyD >= rows || grid[gx][gyD] !== 0) {
                canFall = false;
              }

              break;
            }
          }
          if (!canFall) break;
        }

        if (canFall) {
          for (let r = 0; r < 4; r++)
            for (let c = 0; c < 4; c++)
              if (piece.pattern[r][c] === 1) {
                const gx = x + (c - piece.center.x);
                const gy = y + (r - piece.center.y);
                grid[gx][gy] = 0;
              }

          y++;

          for (let r = 0; r < 4; r++)
            for (let c = 0; c < 4; c++)
              if (piece.pattern[r][c] === 1) {
                const gx = x + (c - piece.center.x);
                const gy = y + (r - piece.center.y);
                grid[gx][gy] = id;
              }

          // Reset delai et nouveau cycle de chute
          fallDelay = Math.max(300, 1000 - (lvl - 1) * 50);
          startFallTimer();
          //Si piece bloquée
        } else {
          // Double for de reecriture inutile
          for (let r = 0; r < 4; r++)
            for (let c = 0; c < 4; c++)
              if (piece.pattern[r][c] === 1) {
                const gx = x + (c - piece.center.x);
                const gy = y + (r - piece.center.y);
                grid[gx][gy] = id;
              }

          let cleared = 0;

          // Effacement des lignes pleines et incrementation du score
          for (let rr = rows - 1; rr >= 0; rr--) {
            let full = true;

            for (let cc = 0; cc < cols; cc++)
              if (grid[cc][rr] === 0) {
                full = false;
                break;
              }

            if (full) {
              cleared++;

              for (let tt = rr; tt > 0; tt--)
                for (let cc = 0; cc < cols; cc++)
                  grid[cc][tt] = grid[cc][tt - 1];

              for (let cc = 0; cc < cols; cc++) grid[cc][0] = 0;

              rr++;
            }
          }

          if (cleared === 1) score += 100;
          if (cleared === 2) score += 300;
          if (cleared === 3) score += 500;
          if (cleared === 4) score += 800;

          lvl = Math.floor(Math.sqrt(score / 200)) + 1;

          let blocked = false;

          // Verification possibilité placement nouvelle pièce
          for (let r = 0; r < 4 && !blocked; r++) {
            for (let c = 0; c < 4 && !blocked; c++) {
              if (piece.pattern[r][c] === 1) {
                const gx = spawnX + (c - piece.center.x);
                const gy = spawnY + (r - piece.center.y);

                if (
                  gx < 0 ||
                  gx >= cols ||
                  gy < 0 ||
                  gy >= rows ||
                  grid[gx][gy] !== 0
                ) {
                  blocked = true;
                }
              }
            }
          }

          if (blocked) {
            // Game over
            action = "wait";
            playTimerRunning = false;
            gameOverBox.style.visibility = "visible";
          } else {
            // Placement nouvelle piece et nouveau cyclage de la chute
            x = spawnX;
            y = spawnY;

            id = nextId;
            piece = structuredClone(
              tetrisPieces[Object.keys(tetrisPieces)[id - 1]]
            );
            nextId = Math.floor(Math.random() * 7) + 1;

            updateNextPiece();

            for (let r = 0; r < 4; r++)
              for (let c = 0; c < 4; c++)
                if (piece.pattern[r][c] === 1) {
                  const gx = x + (c - piece.center.x);
                  const gy = y + (r - piece.center.y);
                  grid[gx][gy] = id;
                }
          }

          fallDelay = Math.max(300, 1000 - (lvl - 1) * 50);
          startFallTimer();
        }
      }
    }

    //grid[x][y] = id;

    // Mise a jour valeurs de jeu

    lvlValue.textContent = lvl;
    scoreValue.textContent = score;
  }

  // Mise a jour grille html et nouveau cycle de processus

  updateGrid();
  requestAnimationFrame(processus);
}

requestAnimationFrame(processus);
