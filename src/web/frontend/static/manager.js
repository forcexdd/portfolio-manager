import * as cookies from "./cookies.mjs";
import * as constants from "./const.mjs";

let deleteButtonHTML = document.getElementById("delete-portfolio");

deleteButtonHTML.onclick = async (e) => {
    e.preventDefault();
    if(!confirm('Are you sure?')) {
        return
    }
    
    let portfolio = cookies.getCookie(constants.portfolioNameCookie);
    
    e.preventDefault()

    let formData = new FormData();
    formData.append(constants.portfolioNameFormKey, portfolio);
    
    try {
        let response = await fetch("/remove_portfolio", {
            method: "POST",
            body: formData
        });


        if (response.ok) {
            document.cookie = `${constants.portfolioNameCookie}=`;
            location.reload()
        } else if (response.status === 409) {
            alert("Error! This name is already taken!");
        }
    } catch (error) {
        console.error("Error submitting form:", error);
    }
}