let addStockButtonHTML = document.getElementById("addStockButton");
let stockNameInputHTML = document.getElementById("selectStocks");
let quantityHTML = document.getElementById("quantityInput");
let chosenListTableTbodyHTML = document.getElementById("chosenListTbody");

function renderStocks(stocks) {
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
    let responseText = "Error! Try again later!";
    let successDivHTML = document.getElementById("success_text");
    successDivHTML.style.color = 'red'
    
    e.preventDefault();
    if (!validPortfolioName(portfolioNameHTML.value)) {
        responseText = "Error! Invalid name!";
        successDivHTML.innerText = responseText;
        return;
    }

    let stocks = getStocks();
    if (stocks.length === 0) {
        responseText = "Error! Invalid stocks!";
        successDivHTML.innerText = responseText;
        return;
    }

    let formData = new FormData();
    formData.append("portfolioName", portfolioNameHTML.value);
    stocks.forEach(stock => {
        formData.append("stocks[]", JSON.stringify(stock));
    });

    try {
        let response = await fetch("/add_portfolio", {
            method: "POST",
            body: formData
        });

        
        if (response.ok) {
            responseText = "Success! Portfolio was added.";
            successDivHTML.style.color = 'green'
        } else if (response.status === 409) {
            responseText = "Error! This name is already taken!";
        }
        
        successDivHTML.innerText = responseText;
        successDivHTML.style.color = 'red'
    } catch (error) {
        console.error("Error submitting form:", error);
    }
};