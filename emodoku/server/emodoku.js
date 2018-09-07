
// We need to store the value that the user has specifically assigned
// to a sudoku cell separately from the value that's displayed for the
// cell since the value that's displayed for the cell might have been
// concluded by the constraint engine.
var givens = [
    "---------",
    "---------",
    "---------",
    "---------",
    "---------",
    "---------",
    "---------",
    "---------",
    "---------",
];

var solutionResponse = null;

// getGiven returns the character friom the givens array.
function getGiven(row, col) {
  return givens[row - 1][col - 1];
}

// setGiven converts value (an integer from 1 through 9) to the corresponding
// digit and stores it in givens.
function setGiven(row, col, value) {
  var r = givens[row - 1];
  r = r.slice(0, col - 1) + "0123456789"[value] + r.slice(col, r.length);
  givens[row - 1] = r;
}

function clearGiven(row, col) {
  var r = givens[row - 1];
  r = r.slice(0, col - 1) + "-" + r.slice(col, r.length);
  givens[row - 1] = r;
}

// row and col are 1 based.
function cellPossibilities(row, col) {
  return solutionResponse.Possibilities[row - 1][col - 1];
}

function makeValueGlyphId(value) {
  return "val" + value;
}

function setupSymbolInputs() {
  var parent = document.getElementById("symbols");
  for (var i = 1; i <= 9; i++) {
    var e = document.createElement("input");
    e.setAttribute("oninput", "changeGlyph(this)");
    e.setAttribute("type", "text");
    e.setAttribute("size", "1");
    e.setAttribute("id", makeValueGlyphId(i));
    e.setAttribute("value", "" + i);
    parent.appendChild(e);
  }
}

// value is an integer from 1 to 9.
function valueToGlyph(value) {
  var e = document.getElementById(makeValueGlyphId(value));
  return e.value;
}

function glyphToValue(glyph) {
  parent = document.getElementById("symbols");
  for (var i = 0; i < parent.childElementCount; i++) {
    var e = parent.children[i];
    if (e.value == glyph) {
      return i + 1;
    }
  }
  return null;
}

// onimput handler for the symbol input elements.
function changeGlyph() {
  updateSudokuGrid()
}

function makeCellId(row, col) {
  return "row" + row + "_" + "col" + col;
}

function makeCellMenuId(row, col) {
  return "MenuFor_" + makeCellId(row, col)
}

function cellIdToRowCol(cellId) {
  var m = cellId.match(/row([1-9])_col([1-9])/);
  return { "row": m[1], "col": m[2] }
}

function setupSudokuGrid() {
  var parent = document.getElementById("sudoku");
  for (var row = 1; row <= 9; row++) {
    var tr = document.createElement("tr");
    parent.appendChild(tr);
    for (var col = 1; col <= 9; col++) {
      var td = document.createElement("td");
      tr.appendChild(td);
      var id = makeCellId(row, col);
      td.setAttribute("id", id);
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
      td.textContent = " ";
      var puSet = document.createElement("popupset");
      td.appendChild(puSet);
      td.setAttribute("onclick", "gridCellPopup(this)");
    }
  }
}

var selectedCell = null;

function gridCellPopup(elt) {
  // It appears that popup dialogs aren't sufficiently standardized.
  clearChooser();
  var rc = cellIdToRowCol(elt.id);
  selectedCell = rc;
  var poss = cellPossibilities(rc.row, rc.col);
  var chooser = document.getElementById("chooser");
  var prose = document.createElement("p");
  chooser.appendChild(prose);
  if (getGiven(rc.row, rc.col) == "-") {
    prose.textContent = "Pick one of these values for the selected cell:";
    for (var i = 0; i < poss.length; i++) {
      var item = document.createElement("span");
      chooser.appendChild(item);
      item.setAttribute("class", "value-choice");
      item.textContent = valueToGlyph(poss[i]);
      item.setAttribute("onclick", "pickCellValue(this)");
    }
  } else {
    prose.textContent = "Do you want to clear the selected cell?";
    var clearButton = document.createElement("div");
    clearButton.setAttribute("class", "clear-button");
    clearButton.textContent = "Clear";
    clearButton.setAttribute("onclick", "clearCellValue()");
    chooser.appendChild(clearButton);
  }
}

function clearChooser() {
  var chooser = document.getElementById("chooser");
  while (chooser.firstChild) {
    chooser.removeChild(chooser.firstChild);
  }
  selectedCell = null;
}

function pickCellValue(elt) {
  var value = glyphToValue(elt.textContent);
  setGiven(selectedCell.row, selectedCell.col, value);
  sendSolverRequest();
  clearChooser();
}

function clearCellValue() {
  clearGiven(selectedCell.row, selectedCell.col);
  sendSolverRequest();
  clearChooser();
}

function updateSudokuGrid() {
  // console.log("update grid");
  for (var row = 1; row <= 9; row++) {
    for (var col = 1; col <= 9; col++) {
      var poss = cellPossibilities(row, col)
      var td = document.getElementById(makeCellId(row, col));
      if (poss.length == 1) {
        td.textContent = valueToGlyph(poss[0]);
        if (getGiven(row, col) == "-") {
          td.classList.remove("given");
        } else {
          td.classList.add("given");
        }
      } else {
        td.textContent = " ";
        td.classList.remove("given");
      }
    }
  }
}

const socket = new WebSocket("ws://" + window.location.host + "/solver");

socket.addEventListener("open", function(event) {
  sendSolverRequest();
});

socket.addEventListener("message", function(event) {
  solutionResponse = JSON.parse(event.data);
  updateSudokuGrid();
});

function sendSolverRequest() {
  var msg = givens.join("\n") + "\n";
  socket.send(msg);
  // console.log("sent\n", msg);
}

window.onload = function() {
  setupSymbolInputs();
  setupSudokuGrid();
};

