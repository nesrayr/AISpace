// Get the editForm element and the current fact ID to edit
const editArticleForm = document.querySelector('#form-update-article')
const articleToEdit = editArticleForm && editArticleForm.dataset.articleid

// Add an event listener to listen for the form submit
editArticleForm && editArticleForm.addEventListener('submit', (event) => {
    // Prevent the default behaviour of the form element
    event.preventDefault()

    // Convert form data into a JavaScript object
    const formData = Object.fromEntries(new FormData(editArticleForm));

    return fetch(`/article/${articleToEdit}`, {
        // Use the PATCH method
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json'
        },
        // Convert the form's Object data into JSON
        body: JSON.stringify(formData),
    })
        .then(() => document.location.href=`/article/${articleToEdit}`)// Redirect to show
})

// Get deleteButton element and the current fact ID to delete
const deleteButton = document.querySelector('#delete-button')
const articleToDelete = deleteButton && deleteButton.dataset.articleid

// Add event listener to listen for button click
deleteButton && deleteButton.addEventListener('click', () => {
    // We ask the user if they are sure they want to delete the fact
    const result = confirm("Are you sure you want to delete this article?")

    // If the user cancels the prompt, we exit here
    if (!result) return

    // If the user confirms that they want to delete, we send a DELETE request
    // URL uses the current fact's ID
    // Lastly, we redirect to index
    return fetch(`/article/${articleToDelete}`, { method: 'DELETE' })
        .then(() => document.location.href="/")
})

// Get the editForm element and the current fact ID to edit
const editLaboratoryForm = document.querySelector('#form-update-laboratory')
const laboratoryToEdit = editLaboratoryForm && editLaboratoryForm.dataset.laboratoryid

// Add an event listener to listen for the form submit
editLaboratoryForm && editLaboratoryForm.addEventListener('submit', (event) => {
    // Prevent the default behaviour of the form element
    event.preventDefault()

    // Convert form data into a JavaScript object
    const formData = Object.fromEntries(new FormData(editLaboratoryForm));

    return fetch(`/laboratory/${laboratoryToEdit}`, {
        // Use the PATCH method
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json'
        },
        // Convert the form's Object data into JSON
        body: JSON.stringify(formData),
    })
        .then(() => document.location.href=`/laboratory/${laboratoryToEdit}`)// Redirect to show
})


