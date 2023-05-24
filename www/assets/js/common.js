
const MessageType = {
    Info: Symbol("info"),
    Error: Symbol("error")
}

function isEmpty(string) {
    if (string === undefined) {
        return true;
    }
    return typeof string === 'string' && string.length === 0;
}

function convertDate(dateString) {
    if (!isEmpty(dateString)) {
        if (dateString.startsWith("0001-01-01")) {
            return "";
        }
        return dateString.slice(0, 10);
    }
    return "";
}

function validateNumber(field) {
    if (field.value === "") return false;

    return !isNaN(parseFloat(field.value));
}

function showMessage(text, msgType, msgLabel) {
    //set the message
    msgLabel.innerHTML = text;

    const clearMessage = function() {
        msgLabel.innerHTML = "";
    }

    switch (msgType) {
        case MessageType.Info:
            msgLabel.style.color = 'green';
            break;
        case MessageType.Error:
            msgLabel.style.color = 'red';
            break;
        default:
            console.log('message type not defined')
    }
    setTimeout(clearMessage, 4000);
}