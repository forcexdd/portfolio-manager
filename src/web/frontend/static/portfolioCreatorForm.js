import * as validators from './validators.mjs';
import * as constants from "./const.mjs";

let addAssetsButtonHTML = document.getElementById("addAssetButton");
let assetNameInputHTML = document.getElementById("selectAssets");
let quantityHTML = document.getElementById("quantityInput");
let chosenListTableTbodyHTML = document.getElementById("chosenListTbody");

quantityHTML.oninput = (e) => validators.onNumberInput(e, quantityHTML);

function renderAssets(assets) {
    assets.forEach((obj) => {
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
            let array = getAssets().filter((element) => element.name !== name);
            chosenListTableTbodyHTML.querySelector(`[id='${name}']`).remove();
            updateAssets(array);
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

let chosenAssets = [];

function updateAssets(array) {
    chosenAssets = array;
}

function getAssets() {
    return chosenAssets;
}

addAssetsButtonHTML.onclick = (e) => {
    if (assetNameInputHTML.value === "--SELECT--")
        return;

    let obj = {};
    obj["name"] = assetNameInputHTML.value;
    obj["quantity"] =quantityHTML.value;

    let newObject = true;
    chosenAssets.forEach((e) => {
        if (e.name === obj.name) {
            e.quantity = obj.quantity;
            newObject = false;
        }
    })

    if (newObject)
        chosenAssets.push(obj);

    renderAssets(chosenAssets);
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

    let assets = getAssets();
    if (assets.length === 0) {
        responseText = "Error! Invalid assets!";
        successDivHTML.innerText = responseText;
        return;
    }

    let formData = new FormData();
    formData.append(constants.portfolioNameFormKey, portfolioNameHTML.value);
    assets.forEach(assets => {
        formData.append(constants.allAssetsFormKey, JSON.stringify(assets));
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
    } catch (error) {
        console.error("Error submitting form:", error);
    }
};