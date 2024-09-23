export function onNumberInput(e, inputHTML) {
    function removeChars(str) {
        let newStr = "";

        for (let i = 0; i < str.length; i++)
            if (str.charCodeAt(i) >= '0'.charCodeAt(0) && str.charCodeAt(i) <= '9'.charCodeAt(0))
                newStr += str[i];

        return newStr;
    }

    let newValue = removeChars(inputHTML.value);
    if (newValue === "0" || newValue === "")
        newValue++;

    inputHTML.value = newValue;
}

export default class validators {
}