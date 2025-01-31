// import * as validators from "./validators.mjs";
import * as cookies from "./cookies.mjs";
import * as constants from "./const.mjs"

let renderButtonHTML = document.getElementById("renderTable");
// let budgetInputHTML = document.getElementById("budgetInput");
let indexHTML = document.getElementById("indexes");
let divToPlaceTableHTML = document.getElementById("render-following-index-table-box");

// budgetInputHTML.oninput = (e) => validators.onNumberInput(e, budgetInputHTML);

function addAsset(node) {
    let quantityHTML = node.getElementsByClassName("asset-quantity")[0];
    let quantity = parseInt(quantityHTML.textContent);
    if (isNaN(quantity)) {
        quantity = 0;
    } else {
        quantity++;
    }
    quantityHTML.textContent = String(quantity);
}

function removeAsset(node) {
    let quantityHTML = node.getElementsByClassName("asset-quantity")[0];
    let quantity = parseInt(quantityHTML.textContent);
    if (isNaN(quantity) || quantity === 0) {
        return;
    }
    quantity--;
    quantityHTML.textContent = String(quantity);
}

function removeUnusedAssets() {
    let table = document.getElementById("following-index-table");
    for (let i = 1; i < table.rows.length; i++) {
        let suggested_fraction = parseFloat(table.rows[i].cells[4].innerText);

        if (suggested_fraction === 0) {
            table.deleteRow(i);
            i--;
        }
    }
}

function saveChanges() {
    let table = document.getElementById("following-index-table");
    let assets = [];
    for (let i = 1; i < table.rows.length; i++) {
        let row = table.rows[i];
        let name = row.cells[0].innerText;
        let quantity = parseInt(row.cells[1].innerText);
        assets.push({
            "name": name,
            "quantity": quantity
        });
    }
    let formData = new FormData();
    formData.append("assets[]", JSON.stringify(assets));
    formData.append("portfolioName", cookies.getCookie(constants.portfolioNameCookie));
    fetch("/update_portfolio", {
        method: "POST",
        body: formData
    }).then(async response => {
        if (response.ok) {
            alert("Portfolio updated successfully");
        } else {
            alert("Error updating portfolio");
        }
    }).catch(e => console.error(e));
}

function colorDifference() {
    let table = document.getElementById("following-index-table");
    let rows = table.getElementsByTagName("tr");
    for (let i = 1; i < rows.length; i++) {
        let cells = rows[i].getElementsByTagName("td");
        let diff = parseFloat(cells[5].innerText);
        let index_fraction = parseFloat(cells[4].innerText)

        if (index_fraction === 0) {
            rows[i].style.color = "red";
        }

        if (Math.abs(diff) <= 0.05) {
        } else if (diff < 0) {
            rows[i].style.color = "blue";
        } else {
            rows[i].style.color = "red";
        }
    }
}

renderButtonHTML.onclick = (e) => {
    e.preventDefault();

    // let budget = budgetInputHTML.value
    let index = indexHTML.value

    if (index === "select") {
        alert("Please select an index");
        return;
    }

    let formData = new FormData();
    formData.append("index", index);
    // formData.append("budget", budget);
    formData.append(constants.portfolioNameFormKey, cookies.getCookie(constants.portfolioNameCookie));

    fetch("/render_following_index_table", {
        method: "POST",
        body: formData
    }).then(async response => {
        divToPlaceTableHTML.innerHTML = await response.text();
        colorDifference();
        for (let addButton of document.getElementsByClassName("add-asset-button")) {
            addButton.onclick = (_) => addAsset(addButton.parentElement.parentElement);
        }

        for (let removeButton of document.getElementsByClassName("remove-asset-button")) {
            removeButton.onclick = (_) => removeAsset(removeButton.parentElement.parentElement);
        }

        document.getElementById("remove-unused-assets").onclick = (_) => removeUnusedAssets();
        document.getElementById("save-changes").onclick = (_) => saveChanges();
    })
        .catch(e => console.error(e));
}