let optionsHTML = document.getElementById("options");

chooseRightSection(optionsHTML);

optionsHTML.onchange = (e) => {
    let domain = document.location.origin;
    window.location.href = domain + '/' + optionsHTML.value
};

function chooseRightSection(optionsHTML) {
    function getFirstEndpoint()
    {
        let domain = document.location.origin;
        let domainSize = domain.length;
        let endpoint = document.location.href;
        let endpointSize = endpoint.length;
        for (let i = domainSize + 1; i < endpointSize; i++)
        {
            if (endpoint[i] === '/' || endpoint[i] === "?")
                return endpoint.slice(domainSize + 1, i);
        }
        return endpoint.slice(domainSize + 1, endpointSize);
    }

    let endpoint = getFirstEndpoint();
    
    switch (endpoint)
    {
        case "manager":
            optionsHTML.value = "manager";
            break;
        case "following_index":
            optionsHTML.value = "following_index";
            break;
        case "add_portfolio":
            optionsHTML.value = "add_portfolio";
            break;
    }
}

let portfolioSelectionHTML = document.getElementById("portfolios");

chooseRightPortfolio(portfolioSelectionHTML);

portfolioSelectionHTML.onchange = (e) => {
    document.cookie = `current_portfolio=${portfolioSelectionHTML.value}`;
    location.reload();
};

function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) === 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

function chooseRightPortfolio(portfolioHTML) {
    let cookie = getCookie("current_portfolio");
    if (!cookie)
        portfolioHTML.value = "select";
    else
        portfolioHTML.value = cookie;
}