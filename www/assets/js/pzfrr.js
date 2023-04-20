window.addEventListener('DOMContentLoaded', (event) => {
    const page = window.location.pathname.substring(1);
    console.log('DOM fully loaded and parsed, page: ', page);
    switch (page) {
        case "index.html":
            createPzfrrTable(); break;
        case "files.html":
            setFilesForm(); break;
        case "configuration.html":
            setConfigurationForm(); break;
        default:
            break;
    }
});

function savePzfrrReport() {
    const xhr = new XMLHttpRequest();
    const url='http://'+ host + ':' + port + '/api/save_pzfrr';
    xhr.open("POST", url);
    xhr.setRequestHeader("Accept", "application/json");
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            console.log(xhr.status);
            // console.log(xhr.responseText);
            fillPzfrrForm(JSON.parse(xhr.responseText))
        }};

    xhr.send(getJsonFromPzfrrForm());
}

function saveAndPublishPzfrrReport() {
    const xhr = new XMLHttpRequest();
    const url='http://'+ host + ':' + port + '/api/save_pzfrr_publish';
    xhr.open("POST", url);
    xhr.setRequestHeader("Accept", "application/json");
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            console.log(xhr.status);
            // console.log(xhr.responseText);
            fillPzfrrForm(JSON.parse(xhr.responseText))
        }};

    xhr.send(getJsonFromPzfrrForm());
}

function generatePzfrrReport() {
    let dateFrom = document.getElementById("pzfrr_date_from").value;
    let dateTo = document.getElementById("pzfrr_date_to").value;

    console.log("Get PZFRR report within dates: ", dateFrom, dateTo);

    const get = async (url, params) => {
        const response = await fetch(url + '?' + new URLSearchParams(params))
        const respData = await response.json()

        return respData
    }

    // Calling it with then:
    get('http://'+ host + ':' + port + '/api/test_pzfrr', {
        dateFrom: dateFrom,
        dateTo: dateTo,
    }).then(respData => {
        console.log(respData)
        fillPzfrrForm(respData)
    })
}

function getPzfrrReport() {
    let dateFrom = document.getElementById("pzfrr_date_from").value;
    let dateTo = document.getElementById("pzfrr_date_to").value;

    console.log("Get PZFRR report within dates: ", dateFrom, dateTo);

    const get = async (url, params) => {
        const response = await fetch(url + '?' + new URLSearchParams(params))
        const respData = await response.json()

        return respData
    }

    // Calling it with then:
    get('http://'+ host + ':' + port + '/api/get_pzfrr', {
        dateFrom: dateFrom,
        dateTo: dateTo,
    }).then(respData => {
        console.log(respData)
        fillPzfrrForm(respData)
    })
}

function fillPzfrrForm(respData) {
    const myJSON = JSON.stringify(respData);
    console.log("-------")
    console.log(myJSON)
    console.log("-------")

    let data = respData["Data"];
    let forecastedCapacityUp = respData["ForecastedCapacityUp"];
    let forecastedCapacityDown = respData["ForecastedCapacityDown"];

    fillPzfrrData(data)
    fillPzfrrTableValues("table_pzfrr_forecasted_capacity_up", forecastedCapacityUp);
    fillPzfrrTableValues("table_pzfrr_forecasted_capacity_down", forecastedCapacityDown);

    // createPzfrrTable();
}

function fillPzfrrData(data) {
    const author = document.getElementById("pzfrr_author");
    author.value = data["Creator"];
    let rev = document.getElementById("pzfrr_rev");
    rev.value = data["Revision"];

    const created = data["Created"];
    const saved = data["Saved"];
    const published = data["Reported"];
    setPzfrrDates(created, saved, published)
}

function setPzfrrDates(created, saved, published) {
    let pzfrr_created = document.getElementById("pzfrr-created");
    let pzfrr_saved = document.getElementById("pzfrr-saved");
    let pzfrr_published = document.getElementById("pzfrr-published");

    pzfrr_created.textContent = "Data utworzenia: " + convertDate(created);
    pzfrr_saved.textContent = "Data zapisu: " + convertDate(saved);
    pzfrr_published.textContent = "Data publikacji: " + convertDate(published);
}

function fillPzfrrTableValues(row, values) {
    clearPzfrrTableValues(row)

    let tr = document.getElementById(row);
    for (let i in values) {
        const index = values[i]["position"];
        tr.insertCell(index).innerHTML = values[i]["Quantity"];
    }
}

function getJsonFromPzfrrForm() {
    //convert object to json string
    const data = pzfrrDataToJson();
    const forecastedCapacityUp = pzfrrTableValuesToJson("table_pzfrr_forecasted_capacity_up");
    const forecastedCapacityDown = pzfrrTableValuesToJson("table_pzfrr_forecasted_capacity_down");


    const obj = {};
    obj.data = data;
    obj.forecastedCapacityUp = forecastedCapacityUp;
    obj.forecastedCapacityDown = forecastedCapacityDown;

    return JSON.stringify(obj);
}

function pzfrrDataToJson() {
    const author = document.getElementById("pzfrr_author").value;
    // const rev = document.getElementById("pzrr_rev").innerHTML;
    const date_from = document.getElementById("pzfrr_date_from").value;
    const date_to = document.getElementById("pzfrr_date_to").value;

    let data = {};
    data.creator = author;
    data.start = date_from;
    data.end = date_to;

    return data;
}

function pzfrrTableValuesToJson(row) {
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

function clearPzfrrTableValues(row) {
    let tr = document.getElementById(row);

    while (tr.cells.length > 1) {
        tr.deleteCell(tr.cells.length - 1)
    }
}

function createPzfrrTable() {
    let thr = document.getElementById("table_pzfrr_header_row");

    // let tr = tb.insertRow(-1);
    thr.insertCell(1).innerHTML = "1";
    thr.insertCell(2).innerHTML = "2";
    thr.insertCell(3).innerHTML = "3";
    thr.insertCell(4).innerHTML = "4";

    let tr = document.getElementById("table_pzfrr_forecasted_capacity_up");
    tr.insertCell(1).innerHTML = "1075.0";
    tr.insertCell(2).innerHTML = "1075.0";
    tr.insertCell(3).innerHTML = "1075.0";
    tr.insertCell(4).innerHTML = "1075.0";

    tr = document.getElementById("table_pzfrr_forecasted_capacity_down");
    tr.insertCell(1).innerHTML = "600.0";
    tr.insertCell(2).innerHTML = "600.0";
    tr.insertCell(3).innerHTML = "600.0";
    tr.insertCell(4).innerHTML = "600.0";

}
