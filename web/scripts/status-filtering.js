'use strict';

const STORAGE_STATE_ID = 'status-filtering';

/**
 * @param {HTMLElement} filtersComponent
 */
function saveStatusFilteringState(filtersComponent) {
    window.localStorage.setItem(STORAGE_STATE_ID, JSON.stringify({
        collapsed: filtersComponent.classList.contains('collapsed'),
    }));
}

function getStatusFilteringState() {
    return JSON.parse(window.localStorage.getItem(STORAGE_STATE_ID));
}

/**
 * @param {HTMLElement} filtersComponent
 */
function updateStatusFilteringState(filtersComponent) {
    let newState = getStatusFilteringState();
    if (newState.collapsed) {
        filtersComponent.classList.add('collapsed');
    } else {
        filtersComponent.classList.remove('collapsed');
    }
}

document.addEventListener('DOMContentLoaded', () => {
    let filtersComponent = document.getElementById('filters');

    try {
        updateStatusFilteringState(filtersComponent);
    } catch (e) {
        // NOTE: this should only happen when there is no existing status filtering data in local storage.
        saveStatusFilteringState(filtersComponent);
    }

    // expand and collapse the filters menu by clicking on its header.
    document.getElementById('filters-toggle').addEventListener('click', event => {
        filtersComponent.classList.toggle('collapsed');
        saveStatusFilteringState(filtersComponent);
    });

    document.querySelectorAll('label.filter').forEach(filter => {
        filter.addEventListener('click', () => {}); // TODO:
    });

    window.addEventListener('storage', event => {
        if (event.key === STORAGE_STATE_ID) {
            updateStatusFilteringState(filtersComponent)
        }
    })
});
