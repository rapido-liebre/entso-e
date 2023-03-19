



// // Initialize the DataTable
// document.ready(function () {
//     document.getElementById("tableBody").DataTable({
//
//         // Enable paging
//         // of the DataTable
//         paging: true
//     });
// });

// $(document).ready(function () {
//     $('dataTable').DataTable();
// })

const allDateFs = document.getElementById("dateAll");
const todayFs = document.getElementById("dateToday");
const last3dFs = document.getElementById("dateLast3Days");
const customDateFs = document.getElementById("dateCustom");
const dateFrom = document.getElementById("dateFrom");
const dateTo = document.getElementById("dateTo");

let rangeButtons = [];
rangeButtons.push(allDateFs, todayFs, last3dFs, customDateFs)


function setActiveRangeBtn(elem) {
    rangeButtons.forEach(element => element.classList.remove("active"));
    elem.classList.add("active");
}

function allFiles() {
    const today = new Date()
    const lastYear = new Date(today)

    lastYear.setDate(lastYear.getDate() - 365)

    dateFrom.valueAsDate = lastYear;
    dateTo.valueAsDate = today;

    setActiveRangeBtn(allDateFs)
}

function todayFiles() {
    const today = new Date()

    dateFrom.valueAsDate = today;
    dateTo.valueAsDate = today;

    setActiveRangeBtn(todayFs)
}

function last3DaysFiles() {
    const today = new Date()
    const last3days = new Date(today)

    last3days.setDate(last3days.getDate() - 3)

    dateFrom.valueAsDate = last3days;
    dateTo.valueAsDate = today;

    setActiveRangeBtn(last3dFs)
}

function customFromFiles() {

}

function customToFiles() {

}