
var selectedCell = null;

// We need to store the value that the user has specifically assigned
// to a sudoku cell separately from the value that's displayed for the
// cell since the value that's displayed for the cell might have been
// concluded by the constraint engine.
var givens = [
    "2-----459",
    "----7--3-",
    "6-59---1-",
    "3--89---1",
    "--2---9--",
    "1---27--5",
    "-1---93-4",
    "-2--3----",
    "583-----2"

/*
    "---------",
    "---------",
    "---------",
    "---------",
    "---------",
    "---------",
    "---------",
    "---------",
    "---------",
*/
];

var solutionResponse = null;

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

function makeCellId(row, col) {
  return "row" + row + "_" + "col" + col;
}

function makeCellMenuId(row, col) {
  return "MenuFor_" + makeCellId(row, col)
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
      var menu = document.createElement("menupopup");
      puSet.appendChild(menu);
      var menuId = makeCellMenuId(row, col);
      td.setAttribute("popup", menuId);
      menu.setAttribute("id", menuId);
      // ***** also need an item for clear
      for (var v = 1; v <= 9; v++) {
        var item = document.createElement("menuitem");
        menu.appendChild(item);
        item.setAttribute("label", "" + v);
      }
    }
  }
}

function updateSudokuGrid() {
  for (var row = 1; row <= 9; row++) {
    for (var col = 1; col <= 9; col++) {
      var poss = cellPossibilities(row, col)
      if (poss.length == 1) {
        var td = document.getElementById(makeCellId(row, col));
        td.textContent = valueToGlyph(poss[0]);
      }
    }
  }
}

const socket = new WebSocket("ws://" + window.location.host + "/solver");
console.log("socket", socket);

socket.addEventListener("open", function(event) {
  console.log("received socket open", event);
});

socket.addEventListener("message", function(event) {
  solutionResponse = JSON.parse(event.data);
  updateSudokuGrid();
});

function sendSolverRequest() {
  var msg = givens.join("\n") + "\n";
  socket.send(msg);
  console.log("sent", msg);
}

window.onload = function() {
  console.log("javascript onload");
  setupSymbolInputs();
  setupSudokuGrid();
};

