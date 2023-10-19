let scoreboard;
let scoreboard1;
let scoreboard2;
let scoreboardP1;
let scoreboardP2;

let matchBoard;
let matchBoardTbody;
let holesTable;
let matchScoreDiv;
let matchScoreP;

let sthlmPlayerId = "";
let gbgPlayerId = "";

let totalScore = { Team1: 0, Team2: 0 };
let localTotalScore = 0;
let localScore = [0, 0, 0, 0, 0, 0, 0, 0, 0];

// ON LOAD
window.addEventListener("DOMContentLoaded", () => {
  scoreboard = $(".scoreboard");
  scoreboard1 = $(".team1");
  scoreboard2 = $(".team2");
  scoreboardP1 = $(".team1 p");
  scoreboardP2 = $(".team2 p");
  matchBoard = $("#matchBoard");
  matchBoardTbody = $("#matchBoard tbody");
  holesTable = $("#holes");
  matchScoreDiv = $("#matchScore");
  matchScoreP = $("#matchScore p");

  // Update scoreboard
  getTotalScore();

  // Click matchBoard
  matchBoardTbody.on("click", function (e) {
    const cell = e.target.closest("td");
    if (!cell) {
      // Quit, not clicked on a cell
      return;
    }
    const row = cell.parentElement;
    updateMatchBoard(row.rowIndex, cell.cellIndex);
  });

  // for initial transition of scoreboard
  scoreboard.removeClass("unfocusScoreboard");
});

// TRIGGERS
async function selectPlayer() {
  // If best ball was selected
  if (sthlmPlayerId == "best") {
    $("input[id=bestball]:checked").prop("checked", false);
    focusScoreboard(false);
    sthlmPlayerId = "";
    gbgPlayerId = "";
  }

  sthlmPlayerId = $("input[name=sthlmPlayer]:checked").attr("id");
  gbgPlayerId = $("input[name=gbgPlayer]:checked").attr("id");

  if (sthlmPlayerId && gbgPlayerId) {
    await getMatchAndUpdateBoard();
  }
}

async function bestBall() {
  $("input[name=sthlmPlayer]:checked").prop("checked", false);
  $("input[name=gbgPlayer]:checked").prop("checked", false);
  sthlmPlayerId = "best";
  gbgPlayerId = "ball";

  await getMatchAndUpdateBoard();
}

function clearPlayers() {
  $("input[name=sthlmPlayer]:checked").prop("checked", false);
  sthlmPlayerId = "";
  $("input[name=gbgPlayer]:checked").prop("checked", false);
  gbgPlayerId = "";
  $("input[id=bestball]:checked").prop("checked", false);
  focusScoreboard(true);
}

// FUNCTIONS
async function getMatchAndUpdateBoard() {
  // Get player match
  const response = await fetch(
    "/match?" +
      new URLSearchParams({
        player1: sthlmPlayerId,
        player2: gbgPlayerId,
      })
  );
  if (response.status != 200) {
    console.log("Status: " + response.status.toString() + " " + response.statusText);
    return;
  }

  const score = await response.json();
  drawMatchBoard(score);

  // Unfocus scoreboard
  focusScoreboard(false);
}

async function updateMatchBoard(row, column) {
  if (!sthlmPlayerId || !gbgPlayerId) {
    return; // Players not selected
  }

  // Update score
  if (matchBoardTbody[0].rows[row].cells[column].classList.contains("checked")) {
    // Cell already selected, reduce score
    localScore[column] = 0;
  } else {
    localScore[column] = row + 1;
  }

  // Update table selections
  for (let i = 0; i < 3; i++) {
    if (i == row) {
      matchBoardTbody[0].rows[i].cells[column].classList.toggle("checked");
    } else {
      matchBoardTbody[0].rows[i].cells[column].classList.remove("checked");
    }
  }
  localTotalScore = localScore.reduce((tot, val) => {
    addVal = val > 0 ? val - 2 : 0;
    return tot + addVal;
  }, 0);

  updateMatchTotalScore();

  // Post score to server
  const dataToSend = JSON.stringify({ player1: sthlmPlayerId, player2: gbgPlayerId, score: localScore });
  const resp = await fetch("/match/", {
    method: "PUT",
    headers: { "Content-type": "application/json" },
    body: dataToSend,
  });
  if (resp.status != 200) {
    console.log("Status: " + resp.status);
  }

  getTotalScore();
}

function drawMatchBoard(score) {
  // Update local score
  localScore = score;
  localTotalScore = score.reduce((tot, val) => {
    addVal = val > 0 ? val - 2 : 0;
    return tot + addVal;
  }, 0);

  // Update match board
  for (let i = 0; i < 9; i++) {
    matchBoardTbody[0].rows[0].cells[i].classList.remove("checked");
    matchBoardTbody[0].rows[1].cells[i].classList.remove("checked");
    matchBoardTbody[0].rows[2].cells[i].classList.remove("checked");
    if (localScore[i] > 0) {
      matchBoardTbody[0].rows[localScore[i] - 1].cells[i].classList.add("checked");
    }
  }

  updateMatchTotalScore();
  getTotalScore();
}

async function getTotalScore() {
  // Get player match
  const response = await fetch("/totalScore");
  if (response.status != 200) {
    console.log("Status: " + response.status.toString() + " " + response.statusText);
    return;
  }

  const score = await response.json();
  totalScore = score;

  w1 = Math.floor((totalScore.Team1 / (totalScore.Team1 + totalScore.Team2)) * 100);
  w2 = 100 - w1;
  scoreboard1[0].style.width = w1.toString() + "%";
  scoreboard2[0].style.width = w2.toString() + "%";

  scoreboardP1[0].textContent = totalScore.Team1;
  scoreboardP2[0].textContent = totalScore.Team2;
}

function focusScoreboard(bool) {
  if (bool) {
    scoreboard.removeClass("unfocusScoreboard");
    matchBoard.removeClass("showBoard");
    holesTable.removeClass("showBoard");
    matchScoreDiv.removeClass("showBoard");
  } else {
    scoreboard.addClass("unfocusScoreboard");
    matchBoard.addClass("showBoard");
    holesTable.addClass("showBoard");
    matchScoreDiv.addClass("showBoard");
  }
}

function updateMatchTotalScore() {
  // Update total match score
  if (localTotalScore == 0) {
    matchScoreP[0].style.color = "#808080";
    matchScoreP[0].style.textAlign = "center";

    matchScoreP[0].textContent = "A/S";
  } else if (localTotalScore < 0) {
    matchScoreP[0].style.color = "#003c82";
    matchScoreP[0].style.textAlign = "left";

    matchScoreP[0].textContent = -localTotalScore.toString();
  } else {
    matchScoreP[0].style.color = "#c81414";
    matchScoreP[0].style.textAlign = "right";

    matchScoreP[0].textContent = localTotalScore.toString();
  }
}
