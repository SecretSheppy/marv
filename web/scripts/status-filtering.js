'use strict';

const STATUS_FILTERING_STORAGE_STATE_ID = 'status-filtering';

/**
 * @param {HTMLElement} filtersComponent
 * @returns {*}
 */
function getFilterInputs(filtersComponent) {
    return filtersComponent.querySelectorAll('input[type="checkbox"]')
}

/**
 * @param {HTMLElement} filtersComponent
 * @return {boolean[]}
 */
function getFilterStates(filtersComponent) {
    return Array.from(getFilterInputs(filtersComponent)).map(input => input.checked);
}

/**
 * @param {HTMLElement} filtersComponent
 * @param {boolean[]} newSates
 */
function updateFilterStates(filtersComponent, newSates) {
    let inputs = getFilterInputs(filtersComponent);
    for (let i = 0; i < inputs.length; i++) {
        inputs[i].checked = newSates[i];
    }
}

/**
 * @param {HTMLElement} filtersComponent
 */
function saveStatusFilteringState(filtersComponent) {
    window.localStorage.setItem(STATUS_FILTERING_STORAGE_STATE_ID, JSON.stringify({
        collapsed: filtersComponent.classList.contains('collapsed'),
        checked: getFilterStates(filtersComponent),
    }));
}

function getStatusFilteringState() {
    return JSON.parse(window.localStorage.getItem(STATUS_FILTERING_STORAGE_STATE_ID));
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
    updateFilterStates(filtersComponent, newState.checked)
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
        filter.addEventListener('click', () => {
            saveStatusFilteringState(filtersComponent);
            // TODO:
            //  - use querySelectorAll to select all mutants to show
            //  - if any conflict regions in this stage have their raw source showing, hide it.
            // TODO:
            //  - use querySelectorAll to select all mutants to hide
            //  - if any conflict regions in this stage have no mutants showing, show the raw source.
        });
    });

    window.addEventListener('storage', event => {
        if (event.key === STATUS_FILTERING_STORAGE_STATE_ID) {
            updateStatusFilteringState(filtersComponent)
        }
    })
});
