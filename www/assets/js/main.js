

window.addEventListener('DOMContentLoaded', (event) => {
    const page = window.location.pathname.substring(1);
    console.log('DOM fully loaded and parsed, page: ', page);
    switch (page) {
        case "files.html":
            setFilesForm(); break;
        case "configuration.html":
            setConfigurationForm(); break;
        default:
            break;
    }
});

class Config {
    constructor(timeInterval, warningSize, redAlertSize, inputDir, outputDir, port, path) {
        this.timeInterval = timeInterval;
        this.warningSize = warningSize;
        this.redAlertSize = redAlertSize;
        this.inputDir = inputDir;
        this.outputDir = outputDir;
        this.port = port;
        this.path = path;
    }
}

function setConfigurationForm() {
    loadConfig();
    setConfigFields();
}

function loadConfig() {
    // const Http = new XMLHttpRequest();
    const url='http://localhost:3055/api/config';
    // Http.open("GET", url);
    // Http.send();

    fetch(url, {method: 'GET'})
        .then((response) => response.json())
        .then((data) => {
            console.log(data);
            saveToLocalStorage(data);
        });

    // Http.onreadystatechange = (e) => {
    //     console.log(Http.responseText)
    //     const respText = Http.responseText;//"{\"Params\":{\"TimeInterval\":5,\"WarningSize\":10,\"RedAlertSize\":3,\"InputDir\":\"/Users/rapido_liebre/GolandProjects/wams_archiver/tests/input\",\"OutputDir\":\"/Users/rapido_liebre/GolandProjects/wams_archiver/tests/output\",\"Port\":\":3055\"},\"Path\":\"/Users/rapido_liebre/GolandProjects/wams_archiver/.env\"}"
    //     console.log(respText)
    //     const jsonObj = JSON.parse(respText);
    //     saveToLocalStorage(data)
    //
    //     populate form
        // loadFormFieldsFromLocalStorage()
    // }
}

/**
 * @param {any} jsonObj
 */
function saveToLocalStorage(jsonObj) {
    const myJSON = JSON.stringify(jsonObj);
    // console.log(myJSON)
    localStorage.setItem("testJSON", myJSON);
}

/**
 * @param {string} key
 * @returns {string}
 */
function getLocalStorageValue(key) {
    // Retrieving data
    let text = localStorage.getItem("testJSON");
    let obj = JSON.parse(text);
    return obj[key];
}

function setConfigFields() {
    // map[form_field]cfg_field
    let map = new Map();

    map.set("config_path", "Path");   // a string key
    map.set("time_interval", "TimeInterval");
    map.set("port", "Port");
    map.set("input_directory", "InputDir");
    map.set("output_directory", "OutputDir");
    map.set("warning_disk_size", "WarningSize");
    map.set("alert_disk_size", "RedAlertSize");
    map.set("bucket_name", "BucketName");
    map.set("minio_dir", "MinIODir");
    map.set("minio_url", "MinIOUrl");
    map.set("access_key_id", "AccessKeyId");
    map.set("secret_access_key", "SecretAccessKey");

    map.forEach( (value, key, map) => {
        let lsv = getLocalStorageValue(value === "Path"? value : "Params");
        if (typeof lsv === 'object' && lsv !== null && !Array.isArray(lsv)) {
            for (const [k, v] of Object.entries(lsv)) {
                // console.log(`${k}: ${v}`);
                if (k === value) {
                    lsv = v;
                    break;
                }
            }
        }

        document.getElementById(key).value = lsv;
    });
}





function setFilesForm() {
    loadBucketName("bucket_name");
    loadFiles();
}

function loadBucketName(key) {
    const url = 'http://localhost:3055/api/buckets';

    fetch(url, {method: 'GET'})
        .then((response) => response.json())
        .then((data) => {
            // console.log(data);
            console.log(data[0].name);
            document.getElementById(key).innerHTML = "Bucket Name: " + data[0].name;
        });
}

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

    let tb = document.getElementById("tableBody");
    let tr = tb.insertRow(-1);

    tr.insertCell(0).innerHTML = fileInfo.name;
    tr.insertCell(1).innerHTML = fileInfo.lastModified;
    tr.insertCell(2).innerHTML = fileInfo.size;
    tr.insertCell(3).innerHTML = downloadBtn;
}

function downloadFile(filename) {
    console.log("Download filename: ", filename);

    const get = async (url, params) => {
        const response = await fetch(url + '?' + new URLSearchParams(params))
        const data = await response.json()

        return data
    }

    // Calling it with then:
    get('http://localhost:3055/api/download', {
        filename: filename.replace(" ", "_"),
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


function hello() {
    alert("Hello");
}

// export { getConfig };
