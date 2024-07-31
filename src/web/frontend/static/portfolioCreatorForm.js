let addStockButtonHTML = document.getElementById("addStockButton");
let quantityHTML = document.getElementById("quantityInput");
let stockNameHTML = document.getElementById("selectStocks")
let currentChosenTableHTML = document.getElementById("chosenList");

let chosenStocks = [];

function updateList(listOfObjects) {
    console.log(listOfObjects)
    listOfObjects.forEach((obj) => {
        let name = obj.name;
        let quantity = obj.quantity;

        let existingRow = document.getElementById(name);
        if (existingRow)
        {
            let existingQuantityHTML = existingRow.getElementsByClassName("quantity")[0];
            existingQuantityHTML.innerText = quantity;
            return;
        }

        let row = document.createElement("tr");
        row.id = name;
        let removeButton = document.createElement("button");
        removeButton.type = "button";
        removeButton.onclick = (e) => {
            listOfObjects = listOfObjects.filter((element) => element.name !== name)
            updateList(listOfObjects);
            document.removeChild(row);
        }

        let thName = document.createElement("th");
        thName.innerText = name;

        let thQuantity = document.createElement("th");
        thQuantity.innerText = quantity;
        thQuantity.className = "quantity";

        let thButton = document.createElement("th");
        thButton.appendChild(removeButton);

        row.appendChild(thName);
        row.appendChild(thQuantity);
        row.appendChild(removeButton);

        currentChosenTableHTML.appendChild(row);
    });
}

addStockButtonHTML.addEventListener("click", (e) => {
    if (stockNameHTML.value === "--SELECT--")
        return;

    let obj = {};
    obj["name"] = stockNameHTML.value;
    obj["quantity"] =quantityHTML.value;

    chosenStocks.filter((e) => e.name !== obj.name)
    chosenStocks.push(obj);
    updateList(chosenStocks);
})