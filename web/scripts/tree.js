'use strict';

const TREE_STORAGE_STATE_ID = 'tree-state';

function showCurrentlyOpenFile() {
    let currentFileNode = document.querySelector('meta[name="current-file"]')
    let fileNode = document.querySelector(`a[href="${currentFileNode.content}"]`);
    expandUpwardsToRoot(fileNode);
    fileNode.scrollIntoView({ block: 'start', inline: 'nearest' });
    fileNode.focus();
}

/**
 * @param {HTMLElement} element
 */
function expandUpwardsToRoot(element) {
    let directoryWrapper = element.closest('.directory-wrapper');
    if (directoryWrapper == null) {
        return;
    }
    directoryWrapper.classList.remove('collapsed');
    expandUpwardsToRoot(directoryWrapper.parentElement)
}

/**
 * @param {HTMLElement} element
 */
function expandDirectoriesInsideElement(element) {
    element.querySelectorAll('.directory-wrapper').forEach(dw => {
        dw.classList.remove('collapsed');
    });
}

/**
 * @param {HTMLElement} element
 */
function collapseDirectoriesInsideElement(element) {
    element.querySelectorAll('.directory-wrapper').forEach(dw => {
        dw.classList.add('collapsed');
    });
}

/**
 * @param {Event} event
 */
function collapseToggle(event) {
    let directoryWrapper = event.target.closest('.directory-wrapper');
    directoryWrapper.classList.toggle('collapsed');
    if (directoryWrapper.classList.contains('collapsed')) {
       collapseDirectoriesInsideElement(directoryWrapper);
    }
    saveTreeState();
}

/**
 * @returns {boolean[]}
 */
function getCurrentTreeDirectoryStates() {
    return Array.from(document.querySelectorAll('.directory-wrapper'))
        .map(element => element.classList.contains('collapsed'))
}

function saveTreeState() {
    let treeBody = document.getElementById('tree-body');
    window.localStorage.setItem(TREE_STORAGE_STATE_ID, JSON.stringify({
        directoryStates: getCurrentTreeDirectoryStates(),
        scroll: {
            top: treeBody.scrollTop,
            left: treeBody.scrollLeft,
        }
    }))
}

function getTreeState() {
    return JSON.parse(localStorage.getItem(TREE_STORAGE_STATE_ID));
}

/**
 * @param {boolean[]} newState
 */
function updateTreeState() {
    let newState = getTreeState();

    let directoryStates = newState.directoryStates;
    let directoryWrappers = document.querySelectorAll('.directory-wrapper');
    for (let i = 0; i < directoryStates.length; i++) {
        if (directoryStates[i]) {
            directoryWrappers[i].classList.add('collapsed');
        } else {
            directoryWrappers[i].classList.remove('collapsed');
        }
    }

    document.getElementById('tree-body').scrollTo(newState.scroll)
}

document.addEventListener('DOMContentLoaded', () => {
    try {
        updateTreeState();
    } catch (e) {
        // NOTE: should only occur when there is no existing tree state.
        saveTreeState();
    }

    // expands or collapses the directory wrapper associated with the clicked toggle.
    document.querySelectorAll('.collapse-toggle').forEach(element => {
        element.addEventListener('click', collapseToggle);
    });

    // scrolls to and focuses the currently open file in the file tree.
    document.getElementById('tree-crosshair-btn').addEventListener('click', () => {
        showCurrentlyOpenFile();
        saveTreeState();
    });

    // expands all directories in the tree.
    document.getElementById('tree-expand-all-btn').addEventListener('click', () => {
        expandDirectoriesInsideElement(document.body);
        saveTreeState();
    });

    // collapses all directories in the tree.
    document.getElementById('tree-collapse-all-btn').addEventListener('click', () => {
        collapseDirectoriesInsideElement(document.body);
        saveTreeState();
    });

    document.getElementById('tree-body').addEventListener('scroll', () => saveTreeState());

    window.addEventListener('storage', event => {
        if (event.key === TREE_STORAGE_STATE_ID) {
            updateTreeState();
        }
    })
});