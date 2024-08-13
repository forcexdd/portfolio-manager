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
            let array = getStocks().filter((element) => element.name !== name);
            chosenListTableTbodyHTML.querySelector(`[id='${name}']`).remove();
            updateStocks(array);
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

function getStocks() {
    return chosenStocks;
}

addStockButtonHTML.onclick = (e) => {
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
};

let submitButtonHTML = document.getElementById("submitButton");
let portfolioNameHTML = document.getElementById("portfolioNameInput");

function validPortfolioName(string) {
    return string.length > 0;
}

submitButtonHTML.onclick = async (e) => {
    e.preventDefault();
    console.log(portfolioNameHTML.value);
    if (!validPortfolioName(portfolioNameHTML.value)) {
        alert("invalid name");
        return;
    }

    let stocks = getStocks();
    if (stocks.length === 0) {
        alert("invalid stocks");
        return;
    }

    const formData = new FormData();

    formData.append("portfolioName", portfolioNameHTML.value);
    stocks.forEach((obj) => {
        formData.append(obj.name, obj.quantity);
    })

    try {
        const response = await fetch("/add_portfolio", {
            method: "POST",
            body: formData,
        })
    } catch (e) {
        console.error(e);
    }
};