function isEmpty(string) {
    if (string === undefined) {
        return true;
    }
    const isEmpty = typeof string === 'string' && string.length === 0;
    return isEmpty;
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
