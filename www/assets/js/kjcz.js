window.addEventListener('DOMContentLoaded', (event) => {
    const page = window.location.pathname.substring(1);
    console.log('DOM fully loaded and parsed, page: ', page);
    switch (page) {
        case "index.html":
            createKjczTable(); break;
        case "files.html":
            setFilesForm(); break;
        case "configuration.html":
            setConfigurationForm(); break;
        default:
            // ReadConfigFile();
            break;
    }
});

// function loadFiles() {
//     const url = 'http://localhost:3055/api/bucket_objects';
//
//     fetch(url, {method: 'GET'})
//         .then((response) => response.json())
//         .then((data) => {
//             for (const row of data) {
//                 // console.log(row);
//                 addRowDataToTable(row);
//             }
//         });
//
//     let tb = document.getElementById("tableBody");
//     tb.addEventListener('click', function (e) {
//         const cell = e.target.closest('td');
//         if (!cell) {return;} // Quit, not clicked on a cell
//         const rowIdx = cell.parentElement.rowIndex -1;
//         const filename = tb.rows[rowIdx].cells[0].innerHTML;
//         // console.log("*** ", filename, "  --- ", cell.innerHTML, rowIdx, cell.cellIndex);
//         downloadFile(filename)
//     });
// }

function saveKjczReport() {
    const xhr = new XMLHttpRequest();
    const url='http://'+ host + ':' + port + '/api/save_kjcz';
    xhr.open("POST", url);
    xhr.setRequestHeader("Accept", "application/json");
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            console.log(xhr.status);
            // console.log(xhr.responseText);
            fillKjczForm(JSON.parse(xhr.responseText))
        }};

    xhr.send(getJsonFromKjczForm());
}

function saveAndPublishKjczReport() {
    const xhr = new XMLHttpRequest();
    const url='http://'+ host + ':' + port + '/api/save_kjcz_publish';
    xhr.open("POST", url);
    xhr.setRequestHeader("Accept", "application/json");
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            console.log(xhr.status);
            // console.log(xhr.responseText);
            fillKjczForm(JSON.parse(xhr.responseText))
        }};

    xhr.send(getJsonFromKjczForm());
}

function generateKjczReport() {
    let dateFrom = document.getElementById("kjcz_date_from").value;
    let dateTo = document.getElementById("kjcz_date_to").value;

    console.log("Get KJCZ report within dates: ", dateFrom, dateTo);

    const get = async (url, params) => {
        const response = await fetch(url + '?' + new URLSearchParams(params))
        const respData = await response.json()

        return respData
    }

    // Calling it with then:
    get('http://'+ host + ':' + port + '/api/test_kjcz', {
        dateFrom: dateFrom,
        dateTo: dateTo,
    }).then(respData => {
        console.log(respData)
        fillKjczForm(respData)
    })
}

function getKjczReport() {
    let dateFrom = document.getElementById("kjcz_date_from").value;
    let dateTo = document.getElementById("kjcz_date_to").value;

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
    })
}

function fillKjczForm(respData) {
    const myJSON = JSON.stringify(respData);
    console.log("-------")
    console.log(myJSON)
    console.log("-------")

    let data = respData["Data"];
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

    fillKjczData(data)
    fillKjczTableHeaderValues("table_kjcz_header_row", data);
    fillKjczTableValues("table_kjcz_mean_value_row", meanValue);
    fillKjczTableValues("table_kjcz_standard_deviation_row", standardDeviation);
    fillKjczTableValues("table_kjcz_1_percintile_row", percentile1);
    fillKjczTableValues("table_kjcz_5_percintile_row", percentile5);
    fillKjczTableValues("table_kjcz_10_percintile_row", percentile10);
    fillKjczTableValues("table_kjcz_90_percintile_row", percentile90);
    fillKjczTableValues("table_kjcz_95_percintile_row", percentile95);
    fillKjczTableValues("table_kjcz_99_percintile_row", percentile99);
    fillKjczTableValues("table_kjcz_frce_outside_lev1_range_up_row", frceOutsideLevel1RangeUp);
    fillKjczTableValues("table_kjcz_frce_outside_lev1_range_down_row", frceOutsideLevel1RangeDown);
    fillKjczTableValues("table_kjcz_frce_outside_lev2_range_up_row", frceOutsideLevel2RangeUp);
    fillKjczTableValues("table_kjcz_frce_outside_lev2_range_down_row", frceOutsideLevel2RangeDown);
    fillKjczTableValues("table_kjcz_frce_exceeded_60_frr_capacity_up_row", frceExceeded60PercOfFRRCapacityUp);
    fillKjczTableValues("table_kjcz_frce_exceeded_60_frr_capacity_down_row", frceExceeded60PercOfFRRCapacityDown);

    // createKjczTable();
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

    kjcz_created.textContent = "Data utworzenia: " + convertDate(created);
    kjcz_saved.textContent = "Data zapisu: " + convertDate(saved);
    kjcz_published.textContent = "Data publikacji: " + convertDate(published);
}

function fillKjczTableHeaderValues(row, values) {
    clearKjczTableValues(row);

    let tr = document.getElementById(row);
    const yearMonths = values["YearMonths"]

    const setHeaderDate = function (value, index, array) {
        tr.insertCell(index+1).innerHTML = value
    };
    yearMonths.forEach(setHeaderDate);
}

function fillKjczTableValues(row, values) {
    clearKjczTableValues(row);

    let tr = document.getElementById(row);

    for (let i in values) {
        const index = values[i]["position"];
        tr.insertCell(index).innerHTML = values[i]["Quantity"];
    }
}

function clearKjczTableValues(row) {
    let tr = document.getElementById(row);

    while (tr.cells.length > 1) {
        tr.deleteCell(tr.cells.length - 1)
    }
}

function getJsonFromKjczForm() {
    //convert object to json string
    const data = kjczDataToJson();
    const meanValue = kjczTableValuesToJson("table_kjcz_mean_value_row");
    const standardDeviation = kjczTableValuesToJson("table_kjcz_standard_deviation_row");
    const percentile1 = kjczTableValuesToJson("table_kjcz_1_percintile_row");
    const percentile5 = kjczTableValuesToJson("table_kjcz_5_percintile_row");
    const percentile10 = kjczTableValuesToJson("table_kjcz_10_percintile_row");
    const percentile90 = kjczTableValuesToJson("table_kjcz_90_percintile_row");
    const percentile95 = kjczTableValuesToJson("table_kjcz_95_percintile_row");
    const percentile99 = kjczTableValuesToJson("table_kjcz_99_percintile_row");
    const frceOutsideLevel1RangeUp = kjczTableValuesToJson("table_kjcz_frce_outside_lev1_range_up_row");
    const frceOutsideLevel1RangeDown = kjczTableValuesToJson("table_kjcz_frce_outside_lev1_range_down_row");
    const frceOutsideLevel2RangeUp = kjczTableValuesToJson("table_kjcz_frce_outside_lev2_range_up_row");
    const frceOutsideLevel2RangeDown = kjczTableValuesToJson("table_kjcz_frce_outside_lev2_range_down_row");
    const frceExceeded60PercOfFRRCapacityUp = kjczTableValuesToJson("table_kjcz_frce_exceeded_60_frr_capacity_up_row");
    const frceExceeded60PercOfFRRCapacityDown = kjczTableValuesToJson("table_kjcz_frce_exceeded_60_frr_capacity_down_row");

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

    return JSON.stringify(obj);
}

function kjczDataToJson() {
    const author = document.getElementById("kjcz_author").value;
    // const rev = document.getElementById("kjcz_rev").innerHTML;
    const date_from = document.getElementById("kjcz_date_from").value;
    const date_to = document.getElementById("kjcz_date_to").value;

    let data = {};
    data.creator = author;
    data.start = date_from;
    data.end = date_to;

    return data;
}

function kjczTableValuesToJson(row) {
    let tr = document.getElementById(row);
    let array = [];

    for (let i in tr.cells) {
        if (i === 0) continue;
        // console.log(tr.cells[i].innerHTML);
        let obj = {};
        obj.position = parseInt(i);
        obj.quantity = parseFloat(tr.cells[i].innerHTML);

        array[i-1] = obj;
    }

    return array;
}

function createKjczTable() {
    let thr = document.getElementById("table_kjcz_header_row");

    thr.insertCell(1).innerHTML = "2022-10";
    thr.insertCell(2).innerHTML = "2022-11";
    thr.insertCell(3).innerHTML = "2022-12";

    let tr = document.getElementById("table_kjcz_mean_value_row");
    tr.insertCell(1).innerHTML = "3.309";
    tr.insertCell(2).innerHTML = "1.388";
    tr.insertCell(3).innerHTML = "-1.941";

    tr = document.getElementById("table_kjcz_standard_deviation_row");
    tr.insertCell(1).innerHTML = "56.739";
    tr.insertCell(2).innerHTML = "61.257";
    tr.insertCell(3).innerHTML = "58.645";

    tr = document.getElementById("table_kjcz_1_percintile_row");
    tr.insertCell(1).innerHTML = "-132.749";
    tr.insertCell(2).innerHTML = "-154.430";
    tr.insertCell(3).innerHTML = "-162.567";

    tr = document.getElementById("table_kjcz_5_percintile_row");
    tr.insertCell(1).innerHTML = "-132.749";
    tr.insertCell(2).innerHTML = "-154.430";
    tr.insertCell(3).innerHTML = "-162.567";

    tr = document.getElementById("table_kjcz_10_percintile_row");
    tr.insertCell(1).innerHTML = "-132.749";
    tr.insertCell(2).innerHTML = "-154.430";
    tr.insertCell(3).innerHTML = "-162.567";

    tr = document.getElementById("table_kjcz_90_percintile_row");
    tr.insertCell(1).innerHTML = "-132.749";
    tr.insertCell(2).innerHTML = "-154.430";
    tr.insertCell(3).innerHTML = "-162.567";

    tr = document.getElementById("table_kjcz_95_percintile_row");
    tr.insertCell(1).innerHTML = "-132.749";
    tr.insertCell(2).innerHTML = "-154.430";
    tr.insertCell(3).innerHTML = "-162.567";

    tr = document.getElementById("table_kjcz_99_percintile_row");
    tr.insertCell(1).innerHTML = "-132.749";
    tr.insertCell(2).innerHTML = "-154.430";
    tr.insertCell(3).innerHTML = "-162.567";

    tr = document.getElementById("table_kjcz_frce_outside_lev1_range_up_row");
    tr.insertCell(1).innerHTML = "64";
    tr.insertCell(2).innerHTML = "39";
    tr.insertCell(3).innerHTML = "32";

    tr = document.getElementById("table_kjcz_frce_outside_lev1_range_down_row");
    tr.insertCell(1).innerHTML = "28";
    tr.insertCell(2).innerHTML = "50";
    tr.insertCell(3).innerHTML = "51";

    tr = document.getElementById("table_kjcz_frce_outside_lev2_range_up_row");
    tr.insertCell(1).innerHTML = "6";
    tr.insertCell(2).innerHTML = "8";
    tr.insertCell(3).innerHTML = "0";

    tr = document.getElementById("table_kjcz_frce_outside_lev2_range_down_row");
    tr.insertCell(1).innerHTML = "3";
    tr.insertCell(2).innerHTML = "8";
    tr.insertCell(3).innerHTML = "10";

    tr = document.getElementById("table_kjcz_frce_exceeded_60_frr_capacity_up_row");
    tr.insertCell(1).innerHTML = "7";
    tr.insertCell(2).innerHTML = "2";
    tr.insertCell(3).innerHTML = "0";

    tr = document.getElementById("table_kjcz_frce_exceeded_60_frr_capacity_down_row");
    tr.insertCell(1).innerHTML = "1";
    tr.insertCell(2).innerHTML = "3";
    tr.insertCell(3).innerHTML = "6";
}

function hello(page) {
    alert("Hello " + page);
}
