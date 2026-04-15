'use strict';

document.addEventListener('DOMContentLoaded', () => {
    [...document.getElementsByClassName('collapse-toggle')].forEach(element => {
        element.addEventListener('click', () => {
            let classList = element.closest('.directory-wrapper').classList;
            classList.toggle('collapsed');
            classList.toggle('expanded');
        })
    })
})