// import * as validators from "./validators.mjs";
import * as Cookies from "./cookies.mjs";

let renderButtonHTML = document.getElementById("renderTable");
// let budgetInputHTML = document.getElementById("budgetInput");
let indexHTML = document.getElementById("indexes");
let divToPlaceTableHTML = document.getElementById("render-following-index-table-box");

// budgetInputHTML.oninput = (e) => validators.onNumberInput(e, budgetInputHTML);

function colorDifference() {
    let table = document.getElementById("following-index-table");
    let rows = table.getElementsByTagName("tr");
    for (let i = 1; i < rows.length; i++) {
        let cells = rows[i].getElementsByTagName("td");
        let diff = parseFloat(cells[5].innerText);

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
    formData.append("portfolio", Cookies.getCookie("current_portfolio"));

    fetch("/render_following_index_table", {
        method: "POST",
        body: formData
    }).then(async response => {
        divToPlaceTableHTML.innerHTML = await response.text();
        colorDifference();
    })
        .catch(e => console.error(e));
}