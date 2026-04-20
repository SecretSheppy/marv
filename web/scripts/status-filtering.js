'use strict';

const STATUS_FILTERING_STORAGE_STATE_ID = 'status-filtering';

function styleLastMutants() {
    let bodies = document.getElementById('code-table').querySelectorAll('tbody');
    let lastShowingMutant = 0;
    for (let i = 0; i < bodies.length; i++) {
        bodies[i].classList.remove('last');
        if (bodies[i].classList.contains('mutation') && !bodies[i].classList.contains('hidden')) {
            lastShowingMutant = i;
        }
        if (!bodies[i].classList.contains('mutation') && !bodies[i].classList.contains('hidden')) {
            bodies[lastShowingMutant].classList.add('last');
        }
    }
}

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

function changeMutationVisibility(filtersComponent) {
    if (document.querySelector('meta[name="filtering-enabled"]').content !== 'true') {
        styleLastMutants();
        return;
    }
    let codeTable = document.getElementById('code-table');
    let filters = getFilterInputs(filtersComponent);
    let showQuery = 'tbody[data-status]';
    let hideQuery = 'tbody[data-status]';
    for (let i = 0; i < filters.length; i++) {
        let part = `:not([data-status="${filters[i].name}"])`;
        if (filters[i].checked) {
            hideQuery += part;
        } else {
            showQuery += part;
        }
    }
    let conflictIds = new Set();
    codeTable.querySelectorAll(showQuery).forEach(mutant => {
        mutant.classList.remove('hidden');
        conflictIds.add(mutant.getAttribute('data-conflict-id'));
    });
    codeTable.querySelectorAll(hideQuery).forEach(mutant => {
        mutant.classList.add('hidden');
        conflictIds.add(mutant.getAttribute('data-conflict-id'));
    });
    conflictIds.forEach(id => {
        let mutants = codeTable.querySelectorAll(`tbody[data-conflict-id="${id}"]`);
        let rawShouldBeShown = true;
        for (let i = 0; i < mutants.length; i++) {
            if (!mutants[i].classList.contains('hidden')) {
                rawShouldBeShown = false;
            }
        }
        if (rawShouldBeShown) {
            mutants[0].classList.remove('hidden');
        } else {
            mutants[0].classList.add('hidden');
        }
    });
    styleLastMutants();
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
    updateFilterStates(filtersComponent, newState.checked);
    changeMutationVisibility(filtersComponent);
}

document.addEventListener('DOMContentLoaded', () => {
    let filtersComponent = document.getElementById('filters');

    try {
        updateStatusFilteringState(filtersComponent);
    } catch (e) {
        // NOTE: this should only happen when there is no existing status filtering data in local storage.
        styleLastMutants();
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
            changeMutationVisibility(filtersComponent);
        });
    });

    window.addEventListener('storage', event => {
        if (event.key === STATUS_FILTERING_STORAGE_STATE_ID) {
            updateStatusFilteringState(filtersComponent)
        }
    })
});
