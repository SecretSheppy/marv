'use strict';

function collapseToggle(event) {
    let directoryWrapper = event.target.closest('.directory-wrapper');
    directoryWrapper.classList.toggle('collapsed');
    if (directoryWrapper.classList.contains('collapsed')) {
        directoryWrapper.querySelectorAll('.directory-wrapper').forEach(dw => {
            dw.classList.add('collapsed');
        });
    }
}

document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.collapse-toggle').forEach(element => {
        element.addEventListener('click', collapseToggle);
    });
});