'use strict';

function dataValue(element, item) {
    return element.childNodes.item(item).getAttribute('data-value')
}

function shouldReverse(titles, column) {
    let data = titles.childNodes.item(column).getAttribute('data-asc');
    titles.childNodes.forEach(node => node.setAttribute('data-asc', ''))
    if (data == null || data === '' || data === 'true') {
        titles.childNodes.item(column).setAttribute('data-asc', 'false');
        return false;
    }
    titles.childNodes.item(column).setAttribute('data-asc', 'true');
    return true;
}

function sortTableByColumn(event) {
    let table = event.target.closest('table');
    let tbody = table.childNodes.item(0);
    let elements = [...tbody.childNodes];
    let titles = elements.splice(0, 1)[0];

    let column = event.target.getAttribute('data-column');
    let reverse = shouldReverse(titles, column);
    console.log(reverse)
    let sorted = elements.sort((a, b) => {
        if (reverse) {
            return dataValue(a, column) - dataValue(b, column);
        }
        return dataValue(b, column) - dataValue(a, column);
    });

    let cloned = tbody.cloneNode();
    cloned.appendChild(titles)
    for (let i = 0; i < sorted.length; i++) {
        cloned.appendChild(sorted[i])
    }
    table.replaceChild(cloned, tbody);
}

document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('th > div.sortable').forEach(element => {
        element.addEventListener('click', sortTableByColumn)
    })
})