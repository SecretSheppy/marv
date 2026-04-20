'use strict';

/**
 * @type {string}
 */
let lastActiveMutationID = "";

function hideAllEmptyReviews() {
    document.querySelectorAll('.review').forEach(review => {
        if (review.classList.contains('hidden')) {
            return
        }
        let reviewTextarea = review.querySelector('textarea');
        if (reviewTextarea.value === '') {
            review.classList.add('hidden');
        }
    });
}

/**
 * @param {HTMLElement} mutationElement
 */
function showAndFocusReview(mutationElement) {
    lastActiveMutationID = mutationElement.id;
    let reviewElement = mutationElement.querySelector('.review');
    reviewElement.classList.remove('hidden');
    reviewElement.querySelector('textarea').focus();
    mutationElement.scrollIntoView();
}

/**
 * @param {Event} event
 */
function reviewButtonClicked(event) {
    hideAllEmptyReviews();
    showAndFocusReview(event.target.closest('.mutation'));
}

/**
 * @param {Event} event
 */
async function reviewInputBlurEvent(event) {
    let loaderWrapper = event.target.closest('.review').querySelector('.loader-wrapper');
    let loaderText = loaderWrapper.querySelector('.loader-status');
    loaderWrapper.classList.remove('saved');
    loaderText.innerText = 'Saving';
    let mutationId = event.target.closest('.mutation').id;
    let framework = document.querySelector('meta[name="current-framework"]').content;
    let file = document.querySelector('meta[name="current-file"]').content;
    let response = await fetch(`/api/review/${framework}/${mutationId}`, {
        method: 'PUT',
        body: JSON.stringify({
            file: file.replace(`/${framework}/mutants/`, ''),
            review: event.target.value
        }),
    });
    if (response.ok) {
        loaderWrapper.classList.add('saved');
        loaderText.innerText = 'Saved';
    } else {
        alert(`failed to save review "${event.target.value}" for mutation ${mutationId}`);
    }
}

/**
 * @param {KeyboardEvent} event
 */
function keydownEvent(event) {
    if (!['ArrowUp', 'ArrowDown', 'PageUp', 'PageDown'].includes(event.key)) {
        return
    }

    hideAllEmptyReviews();

    if (lastActiveMutationID === "") {
        showAndFocusReview(document.querySelector('.mutation'));
        return
    }

    let mutations = Array.from(document.querySelectorAll('.mutation'));
    let lastActiveIndex = mutations.findIndex(el => el.id === lastActiveMutationID);
    let element = null;

    switch (true) {
        case event.key === 'ArrowUp' || event.key === 'PageUp':
            while (element === null) {
                lastActiveIndex--;
                if (lastActiveIndex < 0) {
                    lastActiveIndex = mutations.length - 1
                }
                if (!mutations[lastActiveIndex].classList.contains('hidden')) {
                    element = mutations[lastActiveIndex];
                }
            }
            break;
        case event.key === 'ArrowDown' || event.key === 'PageDown':
            while (element === null) {
                lastActiveIndex++;
                if (lastActiveIndex >= mutations.length) {
                    lastActiveIndex = 0;
                }
                if (!mutations[lastActiveIndex].classList.contains('hidden')) {
                    element = mutations[lastActiveIndex];
                }
            }
            break;
    }

    showAndFocusReview(element);
}

document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.review-btn').forEach(btn => {
        btn.addEventListener('click', reviewButtonClicked);
    });

    document.querySelectorAll('.review textarea').forEach(e => {
        e.addEventListener('blur', reviewInputBlurEvent);
    });

    document.addEventListener('keydown', keydownEvent);
});