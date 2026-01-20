const cols = 10;
const rows = 20;

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
      [1, 1, 0, 0],
      [0, 1, 1, 0],
      [0, 0, 0, 0],
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

for (let c = 0; c < cols; c++) {
  const col = [];
  for (let r = 0; r < rows; r++) {
    col.push(0);
  }
  grid.push(col);
}

const container = document.getElementById("grid");
const lvlValue = document.getElementById("lvlValue");
const scoreValue = document.getElementById("scoreValue");
const playButton = document.getElementById("playButton");
const pauseButton = document.getElementById("pauseButton");
const pauseRestart = document.getElementById("pauseRestart");
const pauseResume = document.getElementById("pauseResume");

for (let r = 0; r < rows; r++) {
  for (let c = 0; c < cols; c++) {
    const div = document.createElement("div");
    div.className = "cell";
    div.dataset.id = grid[c][r];
    container.appendChild(div);
  }
}

const fpsValue = document.getElementById("fpsValue");
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

function updateGrid() {
  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      const index = r * cols + c;
      container.children[index].dataset.id = grid[c][r];
    }
  }
  fpsCounter++;
}

let command = null;

document.addEventListener("keydown", (event) => {
  if (event.repeat) return;
  if (event.key === "ArrowLeft") command = "left";
  if (event.key === "ArrowRight") command = "right";
  if (event.key === "ArrowUp") command = "up";
  if (event.key === "ArrowDown") command = "down";
  if (event.key === " ") command = "start";
  if (event.key === "Escape") command = "pause";
});

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

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

let fallPromise = null;

function startFallTimer() {
  fallPromise = delay(1000).then(() => true);
}

startFallTimer();

let spawnX = Math.floor(cols / 2);
let spawnY = 1;
let action = "wait";
let previousAction = null;

let x;
let y;
let id;
let piece;

let lvl;
let score;

async function processus() {
  while (true) {
    const tickPromise = delay(10);

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

    if (action === "start") {
      for (let c = 0; c < cols; c++) {
        for (let r = 0; r < rows; r++) {
          grid[c][r] = 0;
        }
      }

      x = spawnX;
      y = spawnY;
      id = Math.floor(Math.random() * 7) + 1;

      lvl = 1;
      score = 0;

      fallPromise = null;
      startFallTimer();

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
      action = "play";
    }

    if (action === "pause") {
      pauseOverlay.style.visibility = "visible";
    } else {
      pauseOverlay.style.visibility = "hidden";
    }

    if (action === "play") {
      if (command === "left") {
        let canMove = true;

        for (let r = 0; r < 4; r++) {
          for (let c = 0; c < 4; c++) {
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              const gxL = gx - 1;

              if (
                gxL < 0 ||
                (gy >= 0 && grid[gxL][gy] !== 0 && grid[gxL][gy] !== id)
              ) {
                canMove = false;
              }

              break;
            }
          }
          if (!canMove) break;
        }

        if (canMove) {
          for (let r = 0; r < 4; r++)
            for (let c = 0; c < 4; c++)
              if (piece.pattern[r][c] === 1) {
                const gx = x + (c - piece.center.x);
                const gy = y + (r - piece.center.y);
                if (gx >= 0 && gx < cols && gy >= 0 && gy < rows)
                  grid[gx][gy] = 0;
              }

          x--;

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

      if (command === "right") {
        let canMove = true;

        for (let r = 0; r < 4; r++) {
          for (let c = 3; c >= 0; c--) {
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              const gxR = gx + 1;

              if (
                gxR >= cols ||
                (gy >= 0 && grid[gxR][gy] !== 0 && grid[gxR][gy] !== id)
              ) {
                canMove = false;
              }

              break;
            }
          }
          if (!canMove) break;
        }

        if (canMove) {
          for (let r = 0; r < 4; r++)
            for (let c = 0; c < 4; c++)
              if (piece.pattern[r][c] === 1) {
                const gx = x + (c - piece.center.x);
                const gy = y + (r - piece.center.y);
                if (gx >= 0 && gx < cols && gy >= 0 && gy < rows)
                  grid[gx][gy] = 0;
              }

          x++;

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

      if (command === "up") {
        let canMove = true;

        for (let c = 0; c < 4; c++) {
          for (let r = 0; r < 4; r++) {
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              const gyU = gy - 1;

              if (
                gyU < 0 ||
                (gx >= 0 && grid[gx][gyU] !== 0 && grid[gx][gyU] !== id)
              ) {
                canMove = false;
              }

              break;
            }
          }
          if (!canMove) break;
        }

        if (canMove) {
          for (let r = 0; r < 4; r++)
            for (let c = 0; c < 4; c++)
              if (piece.pattern[r][c] === 1) {
                const gx = x + (c - piece.center.x);
                const gy = y + (r - piece.center.y);
                if (gx >= 0 && gx < cols && gy >= 0 && gy < rows)
                  grid[gx][gy] = 0;
              }

          y--;

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

        for (let c = 0; c < 4; c++) {
          for (let r = 3; r >= 0; r--) {
            if (piece.pattern[r][c] === 1) {
              const gx = x + (c - piece.center.x);
              const gy = y + (r - piece.center.y);
              const gyD = gy + 1;

              if (
                gyD >= rows ||
                (gx >= 0 && grid[gx][gyD] !== 0 && grid[gx][gyD] !== id)
              ) {
                canMove = false;
              }

              break;
            }
          }
          if (!canMove) break;
        }

        if (canMove) {
          for (let r = 0; r < 4; r++)
            for (let c = 0; c < 4; c++)
              if (piece.pattern[r][c] === 1) {
                const gx = x + (c - piece.center.x);
                const gy = y + (r - piece.center.y);
                if (gx >= 0 && gx < cols && gy >= 0 && gy < rows)
                  grid[gx][gy] = 0;
              }

          y++;

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

      /*if (fallPromise) {
        const done = await Promise.race([fallPromise, delay(0)]);

        if (done === true) {
          const isBottom = y === rows - 1;
          const hasBlockBelow = !isBottom && grid[x][y + 1] !== 0;

          if (isBottom || hasBlockBelow) {
            grid[x][y] = id;

            for (let r = rows - 1; r >= 0; r--) {
              let full = true;
              for (let c = 0; c < cols; c++) {
                if (grid[c][r] === 0) {
                  full = false;
                  break;
                }
              }

              if (full) {
                for (let rr = r; rr > 0; rr--) {
                  for (let cc = 0; cc < cols; cc++) {
                    grid[cc][rr] = grid[cc][rr - 1];
                  }
                }

                for (let cc = 0; cc < cols; cc++) {
                  grid[cc][0] = 0;
                }

                r++;
              }
            }

            if (grid[spawnX][spawnY] !== 0) {
            } else {
              x = spawnX;
              y = spawnY;
              id = Math.floor(Math.random() * 7) + 1;
              grid[x][y] = id;
            }
          } else {
            grid[x][y] = 0;
            y++;
            grid[x][y] = id;
          }

          startFallTimer();
        }
      }*/

      grid[x][y] = id;

      lvlValue.textContent = lvl;
      scoreValue.textContent = score;
      updateGrid();
    }

    await tickPromise;
  }
}

processus();
