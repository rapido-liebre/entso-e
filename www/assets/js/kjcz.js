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
            createKjczTable(); break;
    }
});


function loadFiles() {
    const url = 'http://localhost:3055/api/bucket_objects';

    fetch(url, {method: 'GET'})
        .then((response) => response.json())
        .then((data) => {
            for (const row of data) {
                // console.log(row);
                addRowDataToTable(row);
            }
        });

    let tb = document.getElementById("tableBody");
    tb.addEventListener('click', function (e) {
        const cell = e.target.closest('td');
        if (!cell) {return;} // Quit, not clicked on a cell
        const rowIdx = cell.parentElement.rowIndex -1;
        const filename = tb.rows[rowIdx].cells[0].innerHTML;
        // console.log("*** ", filename, "  --- ", cell.innerHTML, rowIdx, cell.cellIndex);
        downloadFile(filename)
    });
}

function addRowDataToTable(fileInfo) {
    // Retrieving data
    const downloadBtn = "<button class='btn btn-outline-primary border-0 py-0' type='button' id='" + fileInfo.name + "'><i class='fas fa-download'></i></button>";

    let tb = document.getElementById("kjcz_table1_col");
    
    let tr = tb.insertRow(-1);

    tr.insertCell(0).innerHTML = fileInfo.name;
    tr.insertCell(1).innerHTML = fileInfo.lastModified;
    tr.insertCell(2).innerHTML = fileInfo.size;
    tr.insertCell(3).innerHTML = downloadBtn;
}

function createKjczTable() {
    let thr = document.getElementById("table_kjcz_header_row");

    // let tr = tb.insertRow(-1);
    thr.insertCell(1).innerHTML = "2022-10";
    thr.insertCell(2).innerHTML = "2022-11";
    thr.insertCell(3).innerHTML = "2022-12";
}

function getKjczReport() {
    let dateFrom = document.getElementById("kjcz_date_from").value;
    let dateTo = document.getElementById("kjcz_date_to").value;

    console.log("Get report within dates: ", dateFrom, dateTo);

    const get = async (url, params) => {
        const response = await fetch(url + '?' + new URLSearchParams(params))
        const data = await response.json()

        return data
    }

    // Calling it with then:
    get('http://localhost:3055/api/get_kjcz', {
        dateFrom: dateFrom,
        dateTo: dateTo,
    }).then(data => console.log(data))

    get('http://localhost:3055/api/get_pzrr', {
        dateFrom: dateFrom,
        dateTo: dateTo,
    }).then(data => console.log(data))

    get('http://localhost:3055/api/get_pzfrr', {
        dateFrom: dateFrom,
        dateTo: dateTo,
    }).then(data => console.log(data))

    // const Http = new XMLHttpRequest();
    // const url='http://localhost:3055/api/download';
    // Http.open("GET", url);
    // Http.send();
    //
    // Http.onreadystatechange = (e) => {
    //     console.log(Http.responseText)
    // }
}

function hello(page) {
    alert("Hello " + page);
}