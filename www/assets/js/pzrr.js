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
    const Http = new XMLHttpRequest();
    const url='http://'+ host + ':' + port + '/api/test_pzrr';
    Http.open("GET", url);
    Http.send();

    Http.onreadystatechange = (e) => {
        console.log(Http.responseText)
    }
}

function saveAndPublishPzrrReport() {
    const Http = new XMLHttpRequest();
    const url='http://'+ host + ':' + port + '/api/test_pzrr_publish';
    Http.open("GET", url);
    Http.send();

    Http.onreadystatechange = (e) => {
        console.log(Http.responseText)
    }
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

    createPzrrTable();
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