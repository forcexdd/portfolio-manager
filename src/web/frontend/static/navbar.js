import * as cookies from "./cookies.mjs";
import * as constants from "./const.mjs";

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
    document.cookie = `${constants.portfolioNameCookie}=${portfolioSelectionHTML.value}`;
    location.reload();
};

function chooseRightPortfolio(portfolioHTML) {
    let cookie = cookies.getCookie(constants.portfolioNameCookie);
    if (!cookie)
        portfolioHTML.value = "select";
    else
        portfolioHTML.value = cookie;
}