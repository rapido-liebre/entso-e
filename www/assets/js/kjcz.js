window.addEventListener('DOMContentLoaded', (event) => {
    const page = window.location.pathname.substring(1);
    console.log('DOM fully loaded and parsed, page: ', page);
});

function saveKjczReport() {
    const err = validateKjcz();
    if (err.length > 0) {
        showKjczMessage(err, MessageType.Error)
        return
    }

    const xhr = new XMLHttpRequest();
    const url='http://'+ host + ':' + port + '/api/save_kjcz';
    xhr.open("POST", url);
    xhr.setRequestHeader("Accept", "application/json");
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            console.log(xhr.status);
            // console.log(xhr.responseText);
            if (xhr.status == 200) {
                fillKjczForm(JSON.parse(xhr.responseText))
            }
            else {
                showKjczMessage("Brak komunikacji z serwerem", MessageType.Error);
            }
        }};

    xhr.send(JSON.stringify(getJsonObjectFromKjczForm()));
}

function publishKjczReport() {
    const err = validateKjcz();
    if (err.length > 0) {
        showKjczMessage(err, MessageType.Error)
        return
    }

    const xhr = new XMLHttpRequest();
    const url='http://'+ host + ':' + port + '/api/save_kjcz_publish';
    xhr.open("POST", url);
    xhr.setRequestHeader("Accept", "application/json");
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            console.log(xhr.status);
            // console.log(xhr.responseText);
            if (xhr.status == 200) {
                fillKjczForm(JSON.parse(xhr.responseText))
            }
            else {
                showKjczMessage("Brak komunikacji z serwerem", MessageType.Error);
            }
        }};

    xhr.send(JSON.stringify(getJsonObjectFromKjczForm()));
}

function exportKjczReport() {
    const err = validateKjcz();
    if (err.length > 0) {
        showKjczMessage(err, MessageType.Error)
        return
    }

    const jsonObj = getJsonObjectFromKjczForm();

    const headers = kjczTableHeadersToJson("kjcz_header_m");
    const meanValue = jsonObj["meanValue"];
    const standardDeviation = jsonObj["standardDeviation"];
    const percentile1 = jsonObj["percentile1"];
    const percentile5 = jsonObj["percentile5"];
    const percentile10 = jsonObj["percentile10"];
    const percentile90 = jsonObj["percentile90"];
    const percentile95 = jsonObj["percentile95"];
    const percentile99 = jsonObj["percentile99"];
    const frceOutsideLevel1RangeUp = jsonObj["frceOutsideLevel1RangeUp"];
    const frceOutsideLevel1RangeDown = jsonObj["frceOutsideLevel1RangeDown"];
    const frceOutsideLevel2RangeUp = jsonObj["frceOutsideLevel2RangeUp"];
    const frceOutsideLevel2RangeDown = jsonObj["frceOutsideLevel2RangeDown"];
    const frceExceeded60PercOfFRRCapacityUp = jsonObj["frceExceeded60PercOfFRRCapacityUp"];
    const frceExceeded60PercOfFRRCapacityDown = jsonObj["frceExceeded60PercOfFRRCapacityDown"];


    const csv = [
        "Article: 185.4 - Values of frequency quality evaluation (Part B)\r\n",
        getDataRows(jsonObj, "data"),
        getItemsRow(headers, "Months"),
        getItemsRow(meanValue, "Mean Value"),
        getItemsRow(standardDeviation, "Standard Deviation"),
        getItemsRow(percentile1, "1 - Percentile"),
        getItemsRow(percentile5, "5 - Percentile"),
        getItemsRow(percentile10, "10 - Percentile"),
        getItemsRow(percentile90, "90 - Percentile"),
        getItemsRow(percentile95, "95 - Percentile"),
        getItemsRow(percentile99, "99 - Percentile"),
        "\r\nNo. of Time Intervals",
        getItemsRow(frceOutsideLevel1RangeUp, "FRCE Outside Level 1 Range Up (positive)"),
        getItemsRow(frceOutsideLevel1RangeDown, "FRCE Outside Level 1 Range Down (negative)"),
        getItemsRow(frceOutsideLevel2RangeUp, "FRCE Outside Level 2 Range Up (positive)"),
        getItemsRow(frceOutsideLevel2RangeDown, "FRCE Outside Level 2 Range Down (negative)"),
        getItemsRow(frceExceeded60PercOfFRRCapacityUp, "FRCE Exceeded 60% of FRR Capacity Up (positive)"),
        getItemsRow(frceExceeded60PercOfFRRCapacityDown, "FRCE Exceeded 60% of FRR Capacity Down (negative)")
    ].join('\r\n');
    console.log(csv);

    let currentDate = new Date().toJSON().slice(0, 10);
    let fileName = `raport_kjcz_${currentDate}.csv`;
    const mimeType = 'text/plain';

    downloadFile(csv, fileName, mimeType);

    showKjczMessage("Raport KJCZ zapisany do pliku " + fileName, MessageType.Info);
}

function getDates() {
    const year = document.getElementById("kjcz_year").value;
    const quarter = document.getElementById("kjcz_quarter").value;
    let dateFrom = "";
    let dateTo = "";
    switch (quarter) {
        case '1':
            dateFrom = `${year}-01`;
            dateTo = `${year}-03`;
            break;
        case '2':
            dateFrom = `${year}-04`;
            dateTo = `${year}-06`;
            break;
        case '3':
            dateFrom = `${year}-07`;
            dateTo = `${year}-09`;
            break;
        case '4':
            dateFrom = `${year}-10`;
            dateTo = `${year}-12`;
            break;
    }
    return [dateFrom, dateTo];
}

function getTimeIntervalsParams() {
    const lev1 = document.getElementById("kjcz_frce_out_level1").value;
    const lev2 = document.getElementById("kjcz_frce_out_level2").value;
    const excCapUp = document.getElementById("kjcz_frce_exc60_cap_up").value;
    const excCapDown = document.getElementById("kjcz_frce_exc60_cap_down").value;

    return [lev1, lev2, excCapUp, excCapDown];
}

function createNewKjczReport() {
    clearKjczTableValues();
    const err = validateExtraParametersKjcz();
    if (err.length > 0) {
        showKjczMessage(err, MessageType.Error)
        return
    }

    const [dateFrom, dateTo] = getDates();
    console.log("Create new KJCZ report within dates: ", dateFrom, dateTo);
    const [lev1, lev2, excCapUp, excCapDown] = getTimeIntervalsParams();

    const get = async (url, params) => {
        const response = await fetch(url + '?' + new URLSearchParams(params))
        const respData = await response.json()

        return respData
    }

    // Calling it with then:
    get('http://'+ host + ':' + port + '/api/test_kjcz', {
        dateFrom: dateFrom,
        dateTo: dateTo,
        level1: lev1,
        level2: lev2,
        excCapacityUp: excCapUp,
        excCapacityDown: excCapDown,
    }).then(respData => {
        console.log(respData)
        // if (xhr.status == 200) {
            fillKjczForm(respData)
        // }
        // else {
        //     showKjczMessage("Brak komunikacji z serwerem", MessageType.Error);
        // }

    }).catch(error => {
        showKjczMessage("Brak komunikacji z serwerem", MessageType.Error);
    });
}

function getKjczReport() {
    clearKjczTableValues();

    const [dateFrom, dateTo] = getDates();
    console.log("Get KJCZ report within dates: ", dateFrom, dateTo);

    const get = async (url, params) => {
        const response = await fetch(url + '?' + new URLSearchParams(params))
        const respData = await response.json()

        return respData
    }

    // Calling it with then:
    get('http://'+ host + ':' + port + '/api/get_kjcz', {
        dateFrom: dateFrom,
        dateTo: dateTo,
    }).then(respData => {
        console.log(respData)
        fillKjczForm(respData)
    }).catch(error => {
        showKjczMessage("Brak komunikacji z serwerem", MessageType.Error);
    });
}

function fillKjczForm(respData) {
    const myJSON = JSON.stringify(respData);
    console.log("-------")
    console.log(myJSON)
    console.log("-------")

    let data = respData["Data"];
    if (data["Creator"] === "" && data["Revision"] === 0 && data["YearMonths"] == null) {
        const errMsg = data["Error"]

        if (errMsg.startsWith("connect to db failed: Cannot find password for a user:")) {
            showKjczMessage("Błąd autoryzacji dostępu do bazy danych", MessageType.Warning);
        }
        else if (errMsg.startsWith("error pinging db: ORA-01109: baza danych nie jest otwarta")) {
            showKjczMessage("Brak komunikacji z bazą danych", MessageType.Warning);
        }
        else if (errMsg.startsWith("error pinging db: ORA-01017: niepoprawna nazwa użytkownika/hasło; odmowa zalogowania")) {
            showKjczMessage("Błąd autoryzacji dostępu do bazy danych, niepoprawna nazwa użytkownika/hasło", MessageType.Warning);
        }
        else {
            showKjczMessage("Brak zapisanego raportu KJCZ dla tego zakresu dat", MessageType.Warning);
        }
    }

    let meanValue = respData["MeanValue"];
    let standardDeviation = respData["StandardDeviation"];
    let percentile1 = respData["Percentile1"];
    let percentile5 = respData["Percentile5"];
    let percentile10 = respData["Percentile10"];
    let percentile90 = respData["Percentile90"];
    let percentile95 = respData["Percentile95"];
    let percentile99 = respData["Percentile99"];
    let frceOutsideLevel1RangeUp = respData["FRCEOutsideLevel1RangeUp"];
    let frceOutsideLevel1RangeDown = respData["FRCEOutsideLevel1RangeDown"];
    let frceOutsideLevel2RangeUp = respData["FRCEOutsideLevel2RangeUp"];
    let frceOutsideLevel2RangeDown = respData["FRCEOutsideLevel2RangeDown"];
    let frceExceeded60PercOfFRRCapacityUp = respData["FRCEExceeded60PercOfFRRCapacityUp"];
    let frceExceeded60PercOfFRRCapacityDown = respData["FRCEExceeded60PercOfFRRCapacityDown"];

    fillKjczData(data);
    fillKjczTableHeaderValues("table_kjcz_header_row", data);
    fillKjczTableValues("kjcz_mean_value_", meanValue);
    fillKjczTableValues("kjcz_st_deviation_", standardDeviation);
    fillKjczTableValues("kjcz_percentile1_", percentile1);
    fillKjczTableValues("kjcz_percentile5_", percentile5);
    fillKjczTableValues("kjcz_percentile10_", percentile10);
    fillKjczTableValues("kjcz_percentile90_", percentile90);
    fillKjczTableValues("kjcz_percentile95_", percentile95);
    fillKjczTableValues("kjcz_percentile99_", percentile99);
    fillKjczTableValues("kjcz_frce_out_level1_up_", frceOutsideLevel1RangeUp);
    fillKjczTableValues("kjcz_frce_out_level1_down_", frceOutsideLevel1RangeDown);
    fillKjczTableValues("kjcz_frce_out_level2_up_", frceOutsideLevel2RangeUp);
    fillKjczTableValues("kjcz_frce_out_level2_down_", frceOutsideLevel2RangeDown);
    fillKjczTableValues("kjcz_frce_exc60_cap_up_", frceExceeded60PercOfFRRCapacityUp);
    fillKjczTableValues("kjcz_frce_exc60_cap_down_", frceExceeded60PercOfFRRCapacityDown);
}

function fillKjczData(data) {
    const author = document.getElementById("kjcz_author");
    author.value = data["Creator"];
    let rev = document.getElementById("kjcz_rev");
    rev.value = data["Revision"];

    const created = data["Created"];
    const saved = data["Saved"];
    const published = data["Reported"];
    setKjczDates(created, saved, published)
}

function setKjczDates(created, saved, published) {
    let kjcz_created = document.getElementById("kjcz-created");
    let kjcz_saved = document.getElementById("kjcz-saved");
    let kjcz_published = document.getElementById("kjcz-published");

    kjcz_created.textContent = "Utworzono: " + convertDate(created);
    kjcz_saved.textContent = "Zapisano: " + convertDate(saved);
    kjcz_published.textContent = "Opublikowano: " + convertDate(published);
}

function fillKjczTableHeaderValues(row, values) {
    const yearMonths = values["YearMonths"]
    if (yearMonths == null) {
        return
    }

    for (let i = 1; i <= 3; i++) {
        document.getElementById("kjcz_header_m" + i).value = yearMonths[i-1];
    }
}

function fillKjczTableValues(field, values) {
    for (let i in values) {
        const index = values[i]["Position"];

        document.getElementById(field + index).value = values[i]["Quantity"];
    }
}

function clearKjczTableValues(row) {
    document.getElementById("kjcz_author").value = "";
    document.getElementById("kjcz_rev").value = "";

    for (let i = 1; i <= 3; i++) {
        document.getElementById("kjcz_header_m" + i).value = "";
        document.getElementById("kjcz_mean_value_" + i).value = "";
        document.getElementById("kjcz_st_deviation_" + i).value = "";
        document.getElementById("kjcz_percentile1_" + i).value = "";
        document.getElementById("kjcz_percentile5_" + i).value = "";
        document.getElementById("kjcz_percentile10_" + i).value = "";
        document.getElementById("kjcz_percentile90_" + i).value = "";
        document.getElementById("kjcz_percentile95_" + i).value = "";
        document.getElementById("kjcz_percentile99_" + i).value = "";
        document.getElementById("kjcz_frce_out_level1_up_" + i).value = "";
        document.getElementById("kjcz_frce_out_level1_down_" + i).value = "";
        document.getElementById("kjcz_frce_out_level2_up_" + i).value = "";
        document.getElementById("kjcz_frce_out_level2_down_" + i).value = "";
        document.getElementById("kjcz_frce_exc60_cap_up_" + i).value = "";
        document.getElementById("kjcz_frce_exc60_cap_down_" + i).value = "";
    }

    document.getElementById("kjcz-created").textContent = "Utworzono: ";
    document.getElementById("kjcz-saved").textContent = "Zapisano: ";
    document.getElementById("kjcz-published").textContent = "Opublikowano: ";
}

function getJsonObjectFromKjczForm() {
    //convert object to json string
    const data = kjczDataToJson();
    const meanValue = kjczTableValuesToJson("kjcz_mean_value_");
    const standardDeviation = kjczTableValuesToJson("kjcz_st_deviation_");
    const percentile1 = kjczTableValuesToJson("kjcz_percentile1_");
    const percentile5 = kjczTableValuesToJson("kjcz_percentile5_");
    const percentile10 = kjczTableValuesToJson("kjcz_percentile10_");
    const percentile90 = kjczTableValuesToJson("kjcz_percentile90_");
    const percentile95 = kjczTableValuesToJson("kjcz_percentile95_");
    const percentile99 = kjczTableValuesToJson("kjcz_percentile99_");
    const frceOutsideLevel1RangeUp = kjczTableValuesToJson("kjcz_frce_out_level1_up_");
    const frceOutsideLevel1RangeDown = kjczTableValuesToJson("kjcz_frce_out_level1_down_");
    const frceOutsideLevel2RangeUp = kjczTableValuesToJson("kjcz_frce_out_level2_up_");
    const frceOutsideLevel2RangeDown = kjczTableValuesToJson("kjcz_frce_out_level2_down_");
    const frceExceeded60PercOfFRRCapacityUp = kjczTableValuesToJson("kjcz_frce_exc60_cap_up_");
    const frceExceeded60PercOfFRRCapacityDown = kjczTableValuesToJson("kjcz_frce_exc60_cap_down_");

    const obj = {};
    obj.data = data;
    obj.meanValue = meanValue;
    obj.standardDeviation = standardDeviation;
    obj.percentile1 = percentile1;
    obj.percentile5 = percentile5;
    obj.percentile10 = percentile10;
    obj.percentile90 = percentile90;
    obj.percentile95 = percentile95;
    obj.percentile99 = percentile99;
    obj.frceOutsideLevel1RangeUp = frceOutsideLevel1RangeUp;
    obj.frceOutsideLevel1RangeDown = frceOutsideLevel1RangeDown;
    obj.frceOutsideLevel2RangeUp = frceOutsideLevel2RangeUp;
    obj.frceOutsideLevel2RangeDown = frceOutsideLevel2RangeDown;
    obj.frceExceeded60PercOfFRRCapacityUp = frceExceeded60PercOfFRRCapacityUp;
    obj.frceExceeded60PercOfFRRCapacityDown = frceExceeded60PercOfFRRCapacityDown;

    return obj;
}

function kjczDataToJson() {
    const author = document.getElementById("kjcz_author").value;
    // const rev = document.getElementById("kjcz_rev").innerHTML;
    const [dateFrom, dateTo] = getDates();

    let data = {};
    data.creator = author;
    data.start = dateFrom;
    data.end = dateTo;

    return data;
}

function kjczTableValuesToJson(field) {
    let array = [];

    for (let i = 1; i <= 3; i++) {
        let obj = {};
        obj.position = i;
        obj.quantity = parseFloat(document.getElementById(field + i).value);
        array[i-1] = obj;
    }

    return array;
}

function kjczTableHeadersToJson(field) {
    let array = [];

    for (let i = 1; i <= 3; i++) {
        let obj = {};
        obj.position = i;
        obj.quantity = document.getElementById(field + i).value;
        array[i-1] = obj;
    }

    return array;
}

function validateKjcz() {
    if (document.getElementById("kjcz_author").value === "") {
        return "Błędna wartość w polu Autor";
    }

    for (let i = 1; i <= 3; i++) {
        if (!validateNumber(document.getElementById("kjcz_mean_value_" + i))) return "Błędna wartość w polu Mean Value, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_st_deviation_" + i))) return "Błędna wartość w polu Standard Deviation, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_percentile1_" + i))) return "Błędna wartość w polu Percentile-1, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_percentile5_" + i))) return "Błędna wartość w polu Percentile-5, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_percentile10_" + i))) return "Błędna wartość w polu Percentile-10, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_percentile90_" + i))) return "Błędna wartość w polu Percentile-90, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_percentile95_" + i))) return "Błędna wartość w polu Percentile-95, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_percentile99_" + i))) return "Błędna wartość w polu Percentile-99, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_frce_out_level1_up_" + i))) return "Błędna wartość w polu FRCE Outside Level 1 Range Up, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_frce_out_level1_down_" + i))) return "Błędna wartość w polu FRCE Outside Level 1 Range Down, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_frce_out_level2_up_" + i))) return "Błędna wartość w polu FRCE Outside Level 2 Range Up, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_frce_out_level2_down_" + i))) return "Błędna wartość w polu FRCE Outside Level 2 Range Down, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_frce_exc60_cap_up_" + i))) return "Błędna wartość w polu FRCE Exceeded 60% of FRR Capacity Up, kolumna " + i;
        if (!validateNumber(document.getElementById("kjcz_frce_exc60_cap_down_" + i))) return "Błędna wartość w polu FRCE Exceeded 60% of FRR Capacity Down, kolumna " + i;
    }
    return "";
}

function validateExtraParametersKjcz() {
    const [lev1, lev2, excCapUp, excCapDown] = getTimeIntervalsParams();

    if (!validateNumberOrEmpty(lev1)) return "Błędna wartość w polu FRCE Outside Level 1";
    if (!validateNumberOrEmpty(lev2)) return "Błędna wartość w polu FRCE Outside Level 2";
    if (!validateNumberOrEmpty(excCapUp)) return "Błędna wartość w polu FRCE Exceeded 60% of FRR Capacity Up";
    if (!validateNumberOrEmpty(excCapDown)) return "Błędna wartość w polu FRCE Exceeded 60% of FRR Capacity Down";

    return "";
}

function showKjczMessage(text, msgType) {
    showMessage(text, msgType, document.getElementById("kjcz_message"))
}

// function getCsvFromKjczForm() {
//     const author = document.getElementById("kjcz_author").value;
//     // const rev = document.getElementById("kjcz_rev").innerHTML;
//     const [dateFrom, dateTo] = getDates();
//
//
//     const meanValue = kjczTableValuesToJson("kjcz_mean_value_");
//     const standardDeviation = kjczTableValuesToJson("kjcz_st_deviation_");
//     const percentile1 = kjczTableValuesToJson("kjcz_percentile1_");
//     const percentile5 = kjczTableValuesToJson("kjcz_percentile5_");
//     const percentile10 = kjczTableValuesToJson("kjcz_percentile10_");
//     const percentile90 = kjczTableValuesToJson("kjcz_percentile90_");
//     const percentile95 = kjczTableValuesToJson("kjcz_percentile95_");
//     const percentile99 = kjczTableValuesToJson("kjcz_percentile99_");
//     const frceOutsideLevel1RangeUp = kjczTableValuesToJson("kjcz_frce_out_level1_up_");
//     const frceOutsideLevel1RangeDown = kjczTableValuesToJson("kjcz_frce_out_level1_down_");
//     const frceOutsideLevel2RangeUp = kjczTableValuesToJson("kjcz_frce_out_level2_up_");
//     const frceOutsideLevel2RangeDown = kjczTableValuesToJson("kjcz_frce_out_level2_down_");
//     const frceExceeded60PercOfFRRCapacityUp = kjczTableValuesToJson("kjcz_frce_exc60_cap_up_");
//     const frceExceeded60PercOfFRRCapacityDown = kjczTableValuesToJson("kjcz_frce_exc60_cap_down_");
//
//     const obj = {};
//     obj.data = data;
//     obj.meanValue = meanValue;
//     obj.standardDeviation = standardDeviation;
//     obj.percentile1 = percentile1;
//     obj.percentile5 = percentile5;
//     obj.percentile10 = percentile10;
//     obj.percentile90 = percentile90;
//     obj.percentile95 = percentile95;
//     obj.percentile99 = percentile99;
//     obj.frceOutsideLevel1RangeUp = frceOutsideLevel1RangeUp;
//     obj.frceOutsideLevel1RangeDown = frceOutsideLevel1RangeDown;
//     obj.frceOutsideLevel2RangeUp = frceOutsideLevel2RangeUp;
//     obj.frceOutsideLevel2RangeDown = frceOutsideLevel2RangeDown;
//     obj.frceExceeded60PercOfFRRCapacityUp = frceExceeded60PercOfFRRCapacityUp;
//     obj.frceExceeded60PercOfFRRCapacityDown = frceExceeded60PercOfFRRCapacityDown;
//
//     return obj;
// }

function hello(page) {
    alert("Hello " + page);
}
