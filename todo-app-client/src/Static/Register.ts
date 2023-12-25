import { addSubmitLoginRegisterFormEvent } from "./Utils"

const registerForm=document.querySelector('.form-container') as HTMLFormElement
const register_spinner=document.querySelector('.container>.spinner-conatiner') as HTMLElement

let registerFunc=async(username:string,password:string,warning_txt:HTMLElement)=>( new Promise<boolean>(resolve=>
    setTimeout(()=>{
        (document.querySelector('.container') as HTMLElement).innerHTML=`
        <p class="celebration-txt">Your account has been created. Please click the text bellow to login</p>
        <a href="/">Login</a>
        `
        return resolve(true)
    },1000)    
))

addSubmitLoginRegisterFormEvent(registerForm,register_spinner,registerFunc)