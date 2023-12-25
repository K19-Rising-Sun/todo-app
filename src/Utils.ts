//Remove warning text when any inputs or textare change in a form
function addRemoveWarningTxtEvent(form: HTMLElement): HTMLElement {
    form.querySelectorAll('input[type="text"],textarea').forEach(element => {
        (element as HTMLElement).addEventListener('input', (e) => {
            (form.querySelector('.warning-txt') as HTMLElement).classList.remove('show')
        })
    })
    return form
}

async function addSubmitLoginRegisterFormEvent(form: HTMLFormElement, form_spinner: HTMLElement, cbSubmitFnc: (username: string, password: string, HTMLWarning_txt: HTMLElement) => Promise<boolean>) {
    form.addEventListener('submit', async (e) => {
        e.preventDefault()
        const username = (form.querySelector('#username') as HTMLInputElement).value
        const password = (form.querySelector('#password') as HTMLInputElement).value
        const warning_txt = (form.querySelector('.warning-txt') as HTMLElement)
        if (!username || !password) {
            warning_txt.textContent = 'Please fill in all the fields'
            return
        }
        try {
            form_spinner.classList.toggle('show')
            await cbSubmitFnc(username, password, warning_txt)
        }
        catch (err) {
            console.log(err)
        }
        finally {
            form_spinner.classList.toggle('show')
        }
    })
}

export { addRemoveWarningTxtEvent, addSubmitLoginRegisterFormEvent }
