let deleteButtonHTML = document.getElementById("delete-portfolio");

deleteButtonHTML.onclick = async (e) => {
    e.preventDefault();
    if(!confirm('Are you sure?')) {
        return
    }
    
    let portfolioSelectionHTML = document.getElementById("portfolios");
    let portfolio = portfolioSelectionHTML.value;
    
    e.preventDefault()

    let formData = new FormData();
    formData.append("portfolioName", portfolio);
    
    try {
        let response = await fetch("/remove_portfolio", {
            method: "POST",
            body: formData
        });


        if (response.ok) {
            document.cookie = `current_portfolio=`;
            location.reload()
        } else if (response.status === 409) {
            alert("Error! This name is already taken!");
        }
    } catch (error) {
        console.error("Error submitting form:", error);
    }
}