'use strict';

const DATA_DESCENDING = 'data-descending'

function dataValue(element, item) {
    let value = element.childNodes.item(item).getAttribute('data-value')
    return value === 'NaN' ? -1 : value;
}

function isAscending(titles, column) {
    let data = titles.childNodes.item(column).getAttribute(DATA_DESCENDING);
    titles.childNodes.forEach(node => node.setAttribute(DATA_DESCENDING, ''))
    if (data == null || data === '' || data === 'false') {
        titles.childNodes.item(column).setAttribute(DATA_DESCENDING, 'true');
        return true;
    }
    titles.childNodes.item(column).setAttribute(DATA_DESCENDING, 'false');
    return false;
}

function sortTableByColumn(table, column) {
    let tbody = table.childNodes.item(0);
    let elements = [...tbody.childNodes];
    let titles = elements.splice(0, 1)[0];
    let ascending = isAscending(titles, column);
    let sorted = elements.sort((a, b) => {
        if (ascending) {
            return dataValue(b, column) - dataValue(a, column);
        }
        return dataValue(a, column) - dataValue(b, column);
    });
    let cloned = tbody.cloneNode();
    cloned.appendChild(titles)
    for (let i = 0; i < sorted.length; i++) {
        cloned.appendChild(sorted[i])
    }
    table.replaceChild(cloned, tbody);
    return ascending
}

function sortTableByColumnOnClick(event) {
    let table = event.target.closest('table');
    let column = event.target.closest('div.sortable').getAttribute('data-column');
    let ascending = sortTableByColumn(table, column)
    history.pushState({}, "", "?column=" + column + "&descending=" + ascending)
}

document.addEventListener('DOMContentLoaded', () => {
    let params = new URLSearchParams(window.location.search)
    let column = params.get('column')
    let ascending = params.get('descending')
    if (column != null && ascending != null) {
        let table = document.querySelector('.generic-table');
        let th = table.childNodes.item(0) // tbody
            .childNodes.item(0) // tr:first-child
            .childNodes.item(column) // th for column: column
        th.setAttribute(DATA_DESCENDING, ascending)
        sortTableByColumn(table, column)
    }

    document.querySelectorAll('th > div.sortable').forEach(element => {
        element.addEventListener('click', sortTableByColumnOnClick)
    })
})