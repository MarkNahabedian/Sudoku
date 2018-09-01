
var selectedCell = null;

function setupSymbolInputs() {
  var parent = document.getElementById("symbols");
  for (var i = 1; i <= 9; i++) {
    var e = document.createElement("input");
    e.setAttribute("type", "text");
    e.setAttribute("size", "1");
    e.setAttribute("id", "val" + i);
    e.setAttribute("value", "" + i);
    parent.appendChild(e);
  }
}

function setupSudokuGrid() {
  var parent = document.getElementById("sudoku");
  for (var row = 1; row <= 9; row++) {
    var tr = document.createElement("tr");
    parent.appendChild(tr);
    for (var col = 1; col <= 9; col++) {
      var td = document.createElement("td");
      tr.appendChild(td);
      td.setAttribute("id", "row" + row + "_" + "col" + col);
      var classes = [];
      switch ((row - 1) % 3) {
        case 0: classes.push("top"); break;
        case 1: classes.push("vmiddle"); break;
        case 2: classes.push("bottom"); break;
      }
      switch ((col - 1) % 3) {
        case 0: classes.push("left"); break;
        case 1: classes.push("hmiddle"); break;
        case 2: classes.push("right"); break;
      }
      td.setAttribute("class", classes.join(" "));
    }
  }
}

window.onload = function() {
  console.log("test javascript loaded");
  setupSymbolInputs();
  setupSudokuGrid();
};

