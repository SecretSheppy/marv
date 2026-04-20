'use strict';

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
 * @param {Event} event
 */
function reviewButtonClicked(event) {
    hideAllEmptyReviews();
    let review = event.target.closest('.mutation').querySelector('.review');
    review.classList.remove('hidden');
    review.querySelector('textarea').focus();
}

/**
 * @param {Event} event
 */
async function textAreaBlur(event) {
    let loaderWrapper = event.target.closest('.review').querySelector('.loader-wrapper');
    let loaderText = loaderWrapper.querySelector('.loader-status');
    loaderWrapper.classList.remove('saved');
    loaderText.innerText = 'Saving';
    let mutationId = event.target.closest('.mutation').id;
    let framework = document.querySelector('meta[name="current-framework"]').content;
    let file = document.querySelector('meta[name="current-file"]').content;
    console.log(framework);
    console.log(mutationId);
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

document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.review-btn').forEach(btn => {
        btn.addEventListener('click', reviewButtonClicked);
    });

    document.querySelectorAll('textarea').forEach(txt => {
        txt.addEventListener('blur', textAreaBlur)
    })
});