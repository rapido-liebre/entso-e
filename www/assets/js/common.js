
const MessageType = {
    Info: Symbol("info"),
    Warning: Symbol("warning"),
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

function validateNumberOrEmpty(value) {
    if (value === "") return true;

    return !isNaN(parseFloat(value));
}

function showMessage(text, msgType, msgLabel) {
    //set the message
    msgLabel.innerHTML = text;

    const clearMessage = function() {
        msgLabel.innerHTML = "";
    }

    switch (msgType) {
        case MessageType.Info:
            msgLabel.style.color = 'ForestGreen';
            break;
        case MessageType.Warning:
            msgLabel.style.color = 'orange';
            break;
        case MessageType.Error:
            msgLabel.style.color = 'red';
            break;
        default:
            console.log('message type not defined')
    }
    setTimeout(clearMessage, 8000);
}

function downloadFile(content, fileName, mimeType) {
    const blob = new Blob([content], { type: mimeType });
    const url = URL.createObjectURL(blob);
    const downloadLink = document.createElement('a');

    downloadLink.href = url;
    downloadLink.download = fileName;

    // Append download link to the DOM and trigger a click to start the download
    document.body.appendChild(downloadLink);
    downloadLink.click();

    // Clean up after the download is complete
    document.body.removeChild(downloadLink);
    URL.revokeObjectURL(url);
}

function JSONtoXML(obj) {
    let xml = '';
    for (let prop in obj) {
        xml += obj[prop] instanceof Array ? '' : '<' + prop + '>';
        if (obj[prop] instanceof Array) {
            for (let array in obj[prop]) {
                xml += '\n<' + prop + '>\n';
                xml += JSONtoXML(new Object(obj[prop][array]));
                xml += '</' + prop + '>';
            }
        } else if (typeof obj[prop] == 'object') {
            xml += JSONtoXML(new Object(obj[prop]));
        } else {
            xml += obj[prop];
        }
        xml += obj[prop] instanceof Array ? '' : '</' + prop + '>\n';
    }
    xml = xml.replace(/<\/?[0-9]{1,}>/g, '');
    return xml;
}

// function jsonToCsv(items) {
//     const aaa = items[0]
//     const header = Object.keys(items[0]);
//
//     const headerString = header.join(',');
//
//     // handle null or undefined values here
//     const replacer = (key, value) => value ?? '';
//
//     const rowItems = items.map((row) =>
//         header
//             .map((fieldName) => JSON.stringify(row[fieldName], replacer))
//             .join(',')
//     );
//
//     // join header and body, and break into separate lines
//     const csv = [headerString, ...rowItems].join('\r\n');
//
//     return csv;
// }

function toUpperCase(text) {
    return text.charAt(0).toUpperCase() + text.slice(1);
}

const getDataRows = function (jsonObj, itemName) {
    let obj = jsonObj[itemName];
    const header = Object.keys(obj);
    const upperCaseHeaders = [];
    header.forEach(element => upperCaseHeaders.push(toUpperCase(element)));
    const headerString = upperCaseHeaders.join(',');

    const values = Object.values(obj);
    const valuesString = values.join(',');

    const emptyLine = "";

    const itemRows = [headerString, valuesString, emptyLine].join('\r\n');
    return itemRows;
}

const getItemsRow = function (jsonObjArray, itemName) {
    const row = [];
    row.push(itemName);
    jsonObjArray.forEach((item) => {
        console.log(item);
        row.push(item["quantity"]);
    });

    let itemsRow = row.join(',');
    return itemsRow;
}
