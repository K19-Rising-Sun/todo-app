//Remove warning text when any inputs or textare change in a form
function addRemoveWarningTxtEvent(form: HTMLElement): HTMLElement {
    form.querySelectorAll('input[type="text"],textarea').forEach(element => {
        (element as HTMLElement).addEventListener('input', (e) => {
            (form.querySelector('.warning-txt') as HTMLElement).classList.remove('show')
        })
    })
    return form
}

export { addRemoveWarningTxtEvent }
