window.addEventListener('DOMContentLoaded', (event) => {
    const page = window.location.pathname.substring(1);
    console.log('DOM fully loaded and parsed, page: ', page);
    // switch (page) {
    //     case "index.html":
    //         createKjczTable(); break;
    //     case "files.html":
    //         setFilesForm(); break;
    //     case "configuration.html":
    //         setConfigurationForm(); break;
    //     default:
    //         // ReadConfigFile();
    //         break;
    // }
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

function publishKjczReport() {
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

function exportKjczReport() {

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

function createNewKjczReport() {
    clearKjczTableValues();

    const [dateFrom, dateTo] = getDates();
    console.log("Create new KJCZ report within dates: ", dateFrom, dateTo);

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

function getJsonFromKjczForm() {
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

    return JSON.stringify(obj);
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

function hello(page) {
    alert("Hello " + page);
}
