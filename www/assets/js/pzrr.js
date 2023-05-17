window.addEventListener('DOMContentLoaded', (event) => {
    const page = window.location.pathname.substring(1);
    console.log('DOM fully loaded and parsed, page: ', page);
    switch (page) {
        case "index.html":
            createPzrrTable(); break;
        case "files.html":
            setFilesForm(); break;
        case "configuration.html":
            setConfigurationForm(); break;
        default:
            break;
    }
});

function savePzrrReport() {
    const xhr = new XMLHttpRequest();
    const url='http://'+ host + ':' + port + '/api/save_pzrr';
    xhr.open("POST", url);
    xhr.setRequestHeader("Accept", "application/json");
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            console.log(xhr.status);
            // console.log(xhr.responseText);
            fillPzrrForm(JSON.parse(xhr.responseText))
        }};

    xhr.send(getJsonFromPzrrForm());
}

function saveAndPublishPzrrReport() {
    const xhr = new XMLHttpRequest();
    const url='http://'+ host + ':' + port + '/api/save_pzrr_publish';
    xhr.open("POST", url);
    xhr.setRequestHeader("Accept", "application/json");
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            console.log(xhr.status);
            // console.log(xhr.responseText);
            fillPzrrForm(JSON.parse(xhr.responseText))
        }};

    xhr.send(getJsonFromPzrrForm());
}

function generatePzrrReport() {
    let dateFrom = document.getElementById("pzrr_date_from").value;
    let dateTo = document.getElementById("pzrr_date_to").value;

    console.log("Get PZRR report within dates: ", dateFrom, dateTo);

    const get = async (url, params) => {
        const response = await fetch(url + '?' + new URLSearchParams(params))
        const respData = await response.json()

        return respData
    }

    // Calling it with then:
    get('http://'+ host + ':' + port + '/api/test_pzrr', {
        dateFrom: dateFrom,
        dateTo: dateTo,
    }).then(respData => {
        console.log(respData)
        fillPzrrForm(respData)
    })
}

function getPzrrReport() {
    let dateFrom = document.getElementById("pzrr_date_from").value;
    let dateTo = document.getElementById("pzrr_date_to").value;

    console.log("Get PZRR report within dates: ", dateFrom, dateTo);

    const get = async (url, params) => {
        const response = await fetch(url + '?' + new URLSearchParams(params))
        const respData = await response.json()

        return respData
    }

    // Calling it with then:
    get('http://'+ host + ':' + port + '/api/get_pzrr', {
        dateFrom: dateFrom,
        dateTo: dateTo,
    }).then(respData => {
        console.log(respData)
        fillPzrrForm(respData)
    })
}

function fillPzrrForm(respData) {
    const myJSON = JSON.stringify(respData);
    console.log("-------")
    console.log(myJSON)
    console.log("-------")

    let data = respData["Data"];
    let forecastedCapacityUp = respData["ForecastedCapacityUp"];
    let forecastedCapacityDown = respData["ForecastedCapacityDown"];

    fillPzrrData(data)
    fillPzrrTableHeaderValues("table_pzrr_header_row", data);
    fillPzrrTableValues("table_pzrr_forecasted_capacity_up", forecastedCapacityUp);
    fillPzrrTableValues("table_pzrr_forecasted_capacity_down", forecastedCapacityDown);

    // createPzrrTable();
}

function fillPzrrData(data) {
    const author = document.getElementById("pzrr_author");
    author.value = data["Creator"];
    let rev = document.getElementById("pzrr_rev");
    rev.value = data["Revision"];

    const created = data["Created"];
    const saved = data["Saved"];
    const published = data["Reported"];
    setPzrrDates(created, saved, published)
}

function setPzrrDates(created, saved, published) {
    let pzrr_created = document.getElementById("pzrr-created");
    let pzrr_saved = document.getElementById("pzrr-saved");
    let pzrr_published = document.getElementById("pzrr-published");

    pzrr_created.textContent = "Data utworzenia: " + convertDate(created);
    pzrr_saved.textContent = "Data zapisu: " + convertDate(saved);
    pzrr_published.textContent = "Data publikacji: " + convertDate(published);
}

function fillPzrrTableHeaderValues(row, values) {
    clearPzrrTableValues(row);

    let tr = document.getElementById(row);
    const yearMonths = values["YearMonths"]

    const setHeaderDate = function (value, index, array) {
        tr.insertCell(index+1).innerHTML = value
    };
    if (yearMonths != null) {
        yearMonths.forEach(setHeaderDate);
    }
}

function fillPzrrTableValues(row, values) {
    clearPzrrTableValues(row)

    let tr = document.getElementById(row);
    for (let i in values) {
        const index = values[i]["position"];
        tr.insertCell(index).innerHTML = values[i]["Quantity"];
    }
}

function getJsonFromPzrrForm() {
    //convert object to json string
    const data = pzrrDataToJson();
    const forecastedCapacityUp = pzrrTableValuesToJson("table_pzrr_forecasted_capacity_up");
    const forecastedCapacityDown = pzrrTableValuesToJson("table_pzrr_forecasted_capacity_down");


    const obj = {};
    obj.data = data;
    obj.forecastedCapacityUp = forecastedCapacityUp;
    obj.forecastedCapacityDown = forecastedCapacityDown;

    return JSON.stringify(obj);
}

function pzrrDataToJson() {
    const author = document.getElementById("pzrr_author").value;
    // const rev = document.getElementById("pzrr_rev").innerHTML;
    const date_from = document.getElementById("pzrr_date_from").value;
    const date_to = document.getElementById("pzrr_date_to").value;

    let data = {};
    data.creator = author;
    data.start = date_from;
    data.end = date_to;

    return data;
}

function pzrrTableValuesToJson(row) {
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

function clearPzrrTableValues(row) {
    let tr = document.getElementById(row);

    while (tr.cells.length > 1) {
        tr.deleteCell(tr.cells.length - 1)
    }
}

function createPzrrTable() {
    let thr = document.getElementById("table_pzrr_header_row");

    // let tr = tb.insertRow(-1);
    thr.insertCell(1).innerHTML = "1";
    thr.insertCell(2).innerHTML = "2";
    thr.insertCell(3).innerHTML = "3";
    thr.insertCell(4).innerHTML = "4";

    let tr = document.getElementById("table_pzrr_forecasted_capacity_up");
    tr.insertCell(1).innerHTML = "1500.0";
    tr.insertCell(2).innerHTML = "1500.0";
    tr.insertCell(3).innerHTML = "1500.0";
    tr.insertCell(4).innerHTML = "1500.0";

    tr = document.getElementById("table_pzrr_forecasted_capacity_down");
    tr.insertCell(1).innerHTML = "0.0";
    tr.insertCell(2).innerHTML = "0.0";
    tr.insertCell(3).innerHTML = "0.0";
    tr.insertCell(4).innerHTML = "0.0";


}
