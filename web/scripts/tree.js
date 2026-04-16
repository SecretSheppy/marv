'use strict';

function showCurrentlyOpenFile() {
    let fileNode = document.querySelector(`a[href="${window.location.pathname}"]`);
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
}

document.addEventListener('DOMContentLoaded', () => {
    // expands or collapses the directory wrapper associated with the clicked toggle.
    document.querySelectorAll('.collapse-toggle').forEach(element => {
        element.addEventListener('click', collapseToggle);
    });

    // scrolls to and focuses the currently open file in the file tree.
    document.getElementById('tree-crosshair-btn').addEventListener('click', () => {
        showCurrentlyOpenFile();
    });

    // expands all directories in the tree.
    document.getElementById('tree-expand-all-btn').addEventListener('click', () => {
        expandDirectoriesInsideElement(document.body);
    });

    // collapses all directories in the tree.
    document.getElementById('tree-collapse-all-btn').addEventListener('click', () => {
        collapseDirectoriesInsideElement(document.body);
    });
});