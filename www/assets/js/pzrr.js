window.addEventListener('DOMContentLoaded', (event) => {
    const page = window.location.pathname.substring(1);
    console.log('DOM fully loaded and parsed, page: ', page);
    // switch (page) {
    //     case "index.html":
    //         createPzrrTable(); break;
    //     case "files.html":
    //         setFilesForm(); break;
    //     case "configuration.html":
    //         setConfigurationForm(); break;
    //     default:
    //         break;
    // }
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
            clearPzrrTableValues()
            fillPzrrForm(JSON.parse(xhr.responseText))
        }};

    xhr.send(getJsonFromPzrrForm());
}

function publishPzrrReport() {
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

function exportPzrrReport() {

}

function createNewPzrrReport() {
    clearPzrrTableValues();
}

function getPzrrReport() {
    clearPzrrTableValues();

    const year = document.getElementById("pzrr_year").value;
    const dateFrom = `${year}-01`;
    const dateTo = `${year}-12`;

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

    pzrr_created.textContent = "Utworzono: " + convertDate(created);
    pzrr_saved.textContent = "Zapisano: " + convertDate(saved);
    pzrr_published.textContent = "Opublikowano: " + convertDate(published);
}

function fillPzrrTableValues(row, values) {
    let flowDirection = pzrrGetFlowDirection(row)

    for (let i in values) {
        const index = values[i]["Position"];
        let cell = document.getElementById("pzrr_forecast_cap_" + flowDirection + "_" + index);
        cell.value = values[i]["Quantity"];
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

    const year = document.getElementById("pzrr_year").value;
    const dateFrom = `${year}-01`;
    const dateTo = `${year}-12`;

    let data = {};
    data.creator = author;
    data.start = dateFrom;
    data.end = dateTo;

    return data;
}

function pzrrGetFlowDirection(row) {
    switch (row) {
        case "table_pzrr_forecasted_capacity_up":
            return "up";
        case "table_pzrr_forecasted_capacity_down":
            return "down";
    }
}

function pzrrTableValuesToJson(row) {
    let array = [];
    let flowDirection = pzrrGetFlowDirection(row)

    for (let i = 1; i <= 4; i++) {
        let cell = document.getElementById("pzrr_forecast_cap_" + flowDirection + "_" + i);
        let obj = {};
        obj.position = parseInt(i);
        obj.quantity = parseFloat(cell.value);
        array[i-1] = obj;
    }

    return array;
}

function clearPzrrTableValues() {
    document.getElementById("pzrr_author").value = "";
    document.getElementById("pzrr_rev").value = "";

    for (let i = 1; i <= 4; i++) {
        document.getElementById("pzrr_forecast_cap_up_" + i).value = "";
        document.getElementById("pzrr_forecast_cap_down_" + i).value = "";
    }

    document.getElementById("pzrr-created").textContent = "Utworzono: ";
    document.getElementById("pzrr-saved").textContent = "Zapisano: ";
    document.getElementById("pzrr-published").textContent = "Opublikowano: ";
}

// function createPzrrTable() {
//     let thr = document.getElementById("table_pzrr_header_row");
//
//     // let tr = tb.insertRow(-1);
//     thr.insertCell(1).innerHTML = "1";
//     thr.insertCell(2).innerHTML = "2";
//     thr.insertCell(3).innerHTML = "3";
//     thr.insertCell(4).innerHTML = "4";
//
//     let tr = document.getElementById("table_pzrr_forecasted_capacity_up");
//     tr.insertCell(1).innerHTML = "1500.0";
//     tr.insertCell(2).innerHTML = "1500.0";
//     tr.insertCell(3).innerHTML = "1500.0";
//     tr.insertCell(4).innerHTML = "1500.0";
//
//     tr = document.getElementById("table_pzrr_forecasted_capacity_down");
//     tr.insertCell(1).innerHTML = "0.0";
//     tr.insertCell(2).innerHTML = "0.0";
//     tr.insertCell(3).innerHTML = "0.0";
//     tr.insertCell(4).innerHTML = "0.0";
//
//
// }
