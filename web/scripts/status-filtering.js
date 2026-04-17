'use strict';



document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('filters-toggle').addEventListener('click', event => {
        event.target.closest('.filters-component').classList.toggle('collapsed');
    });

    document.querySelectorAll('label.filter').forEach(filter => {
        filter.addEventListener('click', () => {}); // TODO:
    });
});