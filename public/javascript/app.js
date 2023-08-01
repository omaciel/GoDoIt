// Add click event to fact to show priority
// noinspection SpellCheckingInspection

(() => {
    const priorityWrapper = document.querySelectorAll('.priority-wrapper');
    const toggleBtns = document.querySelectorAll('.priority-toggle')

    for (const ans of priorityWrapper) {
        ans.style.display = 'none';
    }

    for (const btn of toggleBtns) {
        btn.addEventListener('click', (e) => {
            const priority = e.target.parentElement.nextElementSibling;
            priority.style.display = priority.style.display === 'none' ? 'block' : 'none';
        } );
    }

    const editForm = document.querySelector('#form-update-task')
    const taskToEdit = editForm && editForm.dataset.factid


    editForm && editForm.addEventListener('submit', (event) => {
        event.preventDefault()

        const formData = Object.fromEntries(new FormData(editForm));

        return fetch(`/task/${taskToEdit}`, {
                    // Use the PATCH method
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData),
        })
        .then(() => document.location.href=`/task/${taskToEdit}`)
    })




    const deleteButton = document.querySelector('#delete-button')
    const taskToDelete = deleteButton && deleteButton.dataset.factid

    deleteButton && deleteButton.addEventListener('click', () => {
        const result = confirm("Are you sure you want to delete this task?")
        
        if (!result) return

        return fetch(`/task/${taskToDelete}`, { method: 'DELETE' })
                .then(() => document.location.href="/")
    })
})()