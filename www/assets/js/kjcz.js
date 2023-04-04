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
            // createKjczTable();
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
    const Http = new XMLHttpRequest();
    const url='http://localhost:3055/api/test_kjcz';
    Http.open("GET", url);
    Http.send();

    Http.onreadystatechange = (e) => {
        console.log(Http.responseText)
    }
}

function saveAndPublishKjczReport() {
    const Http = new XMLHttpRequest();
    const url='http://localhost:3055/api/test_kjcz_publish';
    Http.open("GET", url);
    Http.send();

    Http.onreadystatechange = (e) => {
        console.log(Http.responseText)
    }
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
    get('http://localhost:3055/api/get_kjcz', {
        dateFrom: dateFrom,
        dateTo: dateTo,
    }).then(respData => {
        console.log(respData)
        fillKjczForm(respData)
    })

    // get('http://localhost:3055/api/get_pzrr', {
    //     dateFrom: dateFrom,
    //     dateTo: dateTo,
    // }).then(data => console.log(data))
    //
    // get('http://localhost:3055/api/get_pzfrr', {
    //     dateFrom: dateFrom,
    //     dateTo: dateTo,
    // }).then(data => console.log(data))

}

function fillKjczForm(respData) {
    const myJSON = JSON.stringify(respData);
    console.log("-------")
    console.log(myJSON)
    console.log("-------")

    // let data = myJSON["Data"];
    // console.log(data);
    // console.log(typeof data);
    // console.log(Array.isArray(data));
    // console.log(data === null);

    createKjczTable();

    // if (typeof lsv === 'object' && lsv !== null && !Array.isArray(lsv)) {
    //     for (const [k, v] of Object.entries(lsv)) {
    //         // console.log(`${k}: ${v}`);
    //         if (k === value) {
    //             lsv = v;
    //             break;
    //         }
    //     }
    // }
    //
    // document.getElementById(key).value = lsv;
    // localStorage.setItem("testJSON", myJSON);

}

function createKjczTable() {
    // let table = document.getElementById("table_kjcz");
    // // var cols = table.cols;
    //
    // // if (cols.length > 1) {
    //
    //     // Getting the rows in table.
    //     var rows = table.rows;
    //
    //     // Removing the column at index(1).
    //     for (var j = 0; j < rows.length; j++) {
    //
    //         // Deleting the ith cell of each row.
    //         rows[j].deleteCell(1);
    //         rows[j].deleteCell(2);
    //         rows[j].deleteCell(3);
    //     }
    // // }
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