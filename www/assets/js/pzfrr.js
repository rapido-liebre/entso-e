window.addEventListener('DOMContentLoaded', (event) => {
    const page = window.location.pathname.substring(1);
    console.log('DOM fully loaded and parsed, page: ', page);
    // switch (page) {
    //     case "index.html":
    //         createPzfrrTable(); break;
    //     case "files.html":
    //         setFilesForm(); break;
    //     case "configuration.html":
    //         setConfigurationForm(); break;
    //     default:
    //         break;
    // }
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
            if (xhr.status == 200) {
                clearPzfrrTableValues()
                fillPzfrrForm(JSON.parse(xhr.responseText))
            }
        }};

    xhr.send(getJsonFromPzfrrForm());
}

function publishPzfrrReport() {
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

function exportPzfrrReport() {

}

function createNewPzfrrReport() {
    clearPzfrrTableValues();
}

function getPzfrrReport() {
    clearPzfrrTableValues();

    const year = document.getElementById("pzfrr_year").value;
    const dateFrom = `${year}-01`;
    const dateTo = `${year}-12`;

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

    pzfrr_created.textContent = "Utworzono: " + convertDate(created);
    pzfrr_saved.textContent = "Zapisano: " + convertDate(saved);
    pzfrr_published.textContent = "Opublikowano: " + convertDate(published);
}

function fillPzfrrTableValues(row, values) {
    let flowDirection = pzfrrGetFlowDirection(row)

    for (let i in values) {
        const index = values[i]["Position"];
        let cell = document.getElementById("pzfrr_forecast_cap_" + flowDirection + "_" + index);
        cell.value = values[i]["Quantity"];
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
    // const rev = document.getElementById("pzfrr_rev").innerHTML;

    const year = document.getElementById("pzfrr_year").value;
    const dateFrom = `${year}-01`;
    const dateTo = `${year}-12`;

    let data = {};
    data.creator = author;
    data.start = dateFrom;
    data.end = dateTo;

    return data;
}

function pzfrrGetFlowDirection(row) {
    switch (row) {
        case "table_pzfrr_forecasted_capacity_up":
            return "up";
        case "table_pzfrr_forecasted_capacity_down":
            return "down";
    }
}

function pzfrrTableValuesToJson(row) {
    let array = [];
    let flowDirection = pzfrrGetFlowDirection(row)

    for (let i = 1; i <= 4; i++) {
        let cell = document.getElementById("pzfrr_forecast_cap_" + flowDirection + "_" + i);
        let obj = {};
        obj.position = parseInt(i);
        obj.quantity = parseFloat(cell.value);
        array[i-1] = obj;
    }

    return array;
}

function clearPzfrrTableValues() {
    document.getElementById("pzfrr_author").value = "";
    document.getElementById("pzfrr_rev").value = "";

    for (let i = 1; i <= 4; i++) {
        document.getElementById("pzfrr_forecast_cap_up_" + i).value = "";
        document.getElementById("pzfrr_forecast_cap_down_" + i).value = "";
    }

    document.getElementById("pzfrr-created").textContent = "Utworzono: ";
    document.getElementById("pzfrr-saved").textContent = "Zapisano: ";
    document.getElementById("pzfrr-published").textContent = "Opublikowano: ";
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
