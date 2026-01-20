const cols = 10;
const rows = 20;

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

let x;
let y;
let id;

let lvl;
let score;

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

let action = "wait";
let previousAction = null;

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

      x = Math.floor(cols / 2);
      y = 0;
      id = Math.floor(Math.random() * 7) + 1;

      lvl = 1;
      score = 0;

      fallPromise = null;
      startFallTimer();

      grid[x][y] = id;
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
        if (x > 0) {
          grid[x][y] = 0;
          x--;
        }
        command = null;
      }

      if (command === "right") {
        if (x < cols - 1) {
          grid[x][y] = 0;
          x++;
        }
        command = null;
      }

      if (command === "up") {
        if (y > 0) {
          grid[x][y] = 0;
          y--;
        }
        command = null;
      }

      if (command === "down") {
        if (y < rows - 1) {
          grid[x][y] = 0;
          y++;
        }
        command = null;
      }

      if (fallPromise) {
        const done = await Promise.race([fallPromise, delay(0)]);

        if (done === true) {
          if (y < rows - 1) {
            grid[x][y] = 0;
            y++;
          }
          startFallTimer();
        }
      }

      grid[x][y] = id;

      lvlValue.textContent = lvl;
      scoreValue.textContent = score;
      updateGrid();
    }

    await tickPromise;
  }
}

processus();
