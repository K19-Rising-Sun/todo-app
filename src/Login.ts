import { addSubmitLoginRegisterFormEvent } from "./Utils"

const loginForm = document.querySelector('.form-container') as HTMLFormElement
const spinner = document.querySelector('.container>.spinner-conatiner') as HTMLElement

let loginFunc = async (username: string, password: string, warning_txt: HTMLElement) => (new Promise<boolean>(resolve =>
    setTimeout(() => {
        return resolve(true)
    }, 1000)
))

addSubmitLoginRegisterFormEvent(loginForm, spinner, loginFunc)
