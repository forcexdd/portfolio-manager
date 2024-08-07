let addStockButtonHTML = document.getElementById("addStockButton");
let stockNameInputHTML = document.getElementById("selectStocks");
let quantityHTML = document.getElementById("quantityInput");
let chosenListTableTbodyHTML = document.getElementById("chosenListTbody");

function renderStocks(stocks) {
    console.log(stocks);
    stocks.forEach((obj) => {
        let name = obj.name;
        let quantity = obj.quantity;

        let existingRow = chosenListTableTbodyHTML.querySelector(`[id='${name}']`);

        if (existingRow) {
            existingRow.querySelector("[id='quantity']").innerText = quantity;
            return;
        }

        let row = document.createElement("tr");
        row.id = name;
        let removeButton = document.createElement("button");
        removeButton.type = "button";
        removeButton.innerText = "Remove";
        removeButton.onclick = (e) => {
            stocks = stocks.filter((element) => element.name !== name);
            chosenListTableTbodyHTML.querySelector(`[id='${name}']`).remove();
            updateStocks(stocks);
            renderStocks(stocks);
        }

        let thName = document.createElement("th");
        thName.innerText = name;

        let thQuantity = document.createElement("th");
        thQuantity.innerText = quantity;
        thQuantity.id = "quantity";

        let thButton = document.createElement("th");
        thButton.appendChild(removeButton);

        row.appendChild(thName);
        row.appendChild(thQuantity);
        row.appendChild(removeButton);

        chosenListTableTbodyHTML.appendChild(row);
    });
}

let chosenStocks = [];

function updateStocks(array) {
    chosenStocks = array;
}

addStockButtonHTML.addEventListener("click", (_) => {
    if (stockNameInputHTML.value === "--SELECT--")
        return;

    let obj = {};
    obj["name"] = stockNameInputHTML.value;
    obj["quantity"] =quantityHTML.value;

    let newObject = true;
    chosenStocks.forEach((e) => {
        if (e.name === obj.name) {
            e.quantity = obj.quantity;
            newObject = false;
        }
    })

    if (newObject)
        chosenStocks.push(obj);

    renderStocks(chosenStocks);
})