import { addRemoveWarningTxtEvent } from "./Utils"
import HTMLEditToDo from "../component/EditTodo"
import HTMLTodo from "../component/Todo"

interface IToDoList {
    [index: string]: {
        title: string,
        category: string,
        description: string,
        is_done: boolean,
    }
}



let toDoList: IToDoList = {}

let stateFilter = 'All'

interface IToDo {
    title: string,
    category: string,
    description: string
}

function checkState(state: string, is_done: boolean): boolean {
    switch (state) {
        case 'Done':
            return is_done
        case 'Undone':
            return !is_done
        default:
            return true
    }
}

function addShowNewToDoFormContainer(element: HTMLElement): HTMLElement {
    element.addEventListener('click', (e) => {
        (document.querySelector('.add-new-todo-container') as HTMLElement).classList.toggle('show')
    })
    return element
}

function renderToDo() {
    const todo_container = document.querySelector('.todo-container') as HTMLElement
    todo_container.innerHTML = `
    <button class="add-todo">
        <span class="material-symbols-outlined">
            add
        </span>
    </button>
    `;
    let add_todo_btn = (todo_container.querySelector('.add-todo') as HTMLElement)
    add_todo_btn = addShowNewToDoFormContainer(add_todo_btn)
    if (Object.keys(toDoList).length === 0) {
        todo_container.innerHTML = '<p class="empty-txt">You have no todo</p>'
        return
    }
    for (const ToDoId in toDoList) {
        const { title, description, category, is_done } = toDoList[ToDoId]
        if (checkState(stateFilter, is_done)) {
            const HTMLToDo = AddHTMLToDoEvents(HTMLTodo(title, category, description, is_done), ToDoId) as HTMLElement
            todo_container.appendChild(HTMLToDo)
        }
    }
}


//Function to call the api
let EditApi = async (toDoId: string):Promise<void>=>{
    const formData = new FormData()
    const {category,title,description}=toDoList[toDoId]
    formData.append('category',category)
    formData.append('title',title)
    formData.append('description',description)
    const response = await fetch('/edittodo',{
        method: 'POST'
    })
    if(response.status<400){
        const {title,category,description} = await response.json()
        toDoList[toDoId]={...toDoList[toDoId],title,category,description}
        return
    }

    throw new Error("Failed to edit todo")
}

let loadToDoFunc = async ():Promise<void> => {
    const response = await fetch('/todo')
    if(response.status<400){
        const newToDoList = await response.json()
        for( const {id,category,title,description,is_done} of newToDoList){
            toDoList[id]={category,title,description,is_done}
        }
        renderToDo()
        return
    }

    throw new Error("Failed to load todo")
}

let deleteToDoFunc = async(toDoId: string): Promise<void> => {
    const response = await fetch('/todo/'+toDoId,{
        method:'DELETE'
    })
    if(response.status<400){
        delete toDoList[toDoId]
        return
    }
    throw new Error('Can\'t delete todo')
}

let checkToDoFunc = (toDoId: string) => new Promise<boolean>((resolve) =>
    setTimeout(() => {
        toDoList[toDoId].is_done = !toDoList[toDoId].is_done
        resolve(true)
    }, 1000)
)

let addToDoFunc = async (newToDo: IToDo): Promise<void> =>{
    const formData = new FormData()
    const {category,title,description}=newToDo
    formData.append('category',category)
    formData.append('title',title)
    formData.append('description',description)
    const response = await fetch('/todo',{
        method: 'POST',
        body: formData
    })
    if(response.status<400){
        const {id,title,category,description,is_done} = await response.json()
        toDoList[id]={title,category,description,is_done}
        renderToDo()
        return
    }

    throw new Error('Failed to add todo')
} 

let searchToDoFunc = async (title: string, category: string): Promise<void>=>{
    const response = await fetch(`/searchtodo?title=${title}&category=${category}`)
    if(response.status<400){
        const newToDoList = await response.json()
        toDoList={}
        for( const {id,category,title,description,is_done} of newToDoList){
            toDoList[id]={category,title,description,is_done}
        }
        renderToDo()
        return
    }

    throw new Error("Failed to search todo")
}

//Inject events into the HTML Elements

//Switch todo element to an edit todo form
function addChangeToEditEvent(HTMLToDo: HTMLElement, ToDoId: string): HTMLElement {
    (HTMLToDo.querySelector('.edit-icon') as HTMLElement).addEventListener('click', (e) => {
        const { title, category, description } = toDoList[ToDoId]
        HTMLToDo.replaceWith(AddHTMLEditToDoEvents(HTMLEditToDo(title, category, description), ToDoId))
    })
    return HTMLToDo
}

//Add delete todo function
function addDeleteToDoEvent(HTMLToDo: HTMLElement, ToDoId: string): HTMLElement {
    const spinner_container = HTMLToDo.querySelector('.spinner-conatiner') as HTMLElement
    (HTMLToDo.querySelector('.delete-icon') as HTMLElement).addEventListener('click', async (e) => {
        try {
            spinner_container.classList.toggle('show')
            await deleteToDoFunc(ToDoId)
            HTMLToDo.remove()
        }
        catch (err) {
            console.log(err)
        }
        finally {
            spinner_container.classList.toggle('show')
        }
    })
    return HTMLToDo
}

//Add check todo function
function addCheckToDoEvent(HTMLToDo: HTMLElement, ToDoId: string): HTMLElement {
    const spinner_container = HTMLToDo.querySelector('.spinner-conatiner') as HTMLElement
    (HTMLToDo.querySelector('.done-icon') as HTMLElement).addEventListener('click', async (e) => {
        try {
            spinner_container.classList.toggle('show')
            await checkToDoFunc(ToDoId)
            HTMLToDo.classList.toggle('complete')
            if (!checkState(stateFilter, toDoList[ToDoId].is_done)) HTMLToDo.remove()
        }
        catch (err) {
            console.log(err)
        }
        finally {
            spinner_container.classList.toggle('show')
        }
    })
    return HTMLToDo
}

//Add all todo element events
function AddHTMLToDoEvents(HTMLToDo: HTMLElement, ToDoId: string): HTMLElement {
    HTMLToDo = addChangeToEditEvent(HTMLToDo, ToDoId)
    HTMLToDo = addDeleteToDoEvent(HTMLToDo, ToDoId)
    HTMLToDo = addCheckToDoEvent(HTMLToDo, ToDoId)
    return HTMLToDo
}

//Change the todo in toDoList whenerver the input changes
function addEditToDoInputEvent(EditHTMLToDo: HTMLElement, ToDoId: string): HTMLElement {
    (EditHTMLToDo.querySelector('input[name="title"]') as HTMLInputElement).addEventListener('input', (e) => {
        toDoList[ToDoId].title = (e.target as HTMLInputElement).value
    });
    (EditHTMLToDo.querySelector('input[name="category"]') as HTMLInputElement).addEventListener('input', (e) => {
        toDoList[ToDoId].category = (e.target as HTMLInputElement).value
    });
    (EditHTMLToDo.querySelector('textarea[name="description"]') as HTMLInputElement).addEventListener('input', (e) => {
        toDoList[ToDoId].description = (e.target as HTMLInputElement).value
    })
    return EditHTMLToDo
}


//Submit change to a todo
function addEditToDoEvent(EditHTMLToDo: HTMLElement, ToDoId: string): HTMLElement {
    const spinner_container = EditHTMLToDo.querySelector('.spinner-conatiner') as HTMLElement
    EditHTMLToDo.addEventListener('submit', async (e) => {
        e.preventDefault()
        const { title, category, description } = toDoList[ToDoId]
        //Cancel if one of the field is empty
        if (!title || !category || !description) {
            EditHTMLToDo.querySelector('.warning-txt')?.classList.toggle('show')
            return
        }
        try {
            spinner_container.classList.toggle('show')
            await EditApi(ToDoId)
            const { title, category, description, is_done } = toDoList[ToDoId]
            EditHTMLToDo.replaceWith(AddHTMLToDoEvents(HTMLTodo(title, category, description, is_done), ToDoId) as HTMLElement)
        }
        catch (err) {
            console.log(err)
        }
        finally {
            spinner_container.classList.toggle('show')
        }
    })
    return EditHTMLToDo
}


//Add all edit todo form events
function AddHTMLEditToDoEvents(EditHTMLToDo: HTMLElement, ToDoId: string): HTMLElement {
    EditHTMLToDo = addEditToDoEvent(EditHTMLToDo, ToDoId)
    EditHTMLToDo = addEditToDoInputEvent(EditHTMLToDo, ToDoId)
    EditHTMLToDo = addRemoveWarningTxtEvent(EditHTMLToDo)
    return EditHTMLToDo
}

//Create a new todo form event
function addCreateNewToDoFormContainerEvent(NewToDoFormContainer: HTMLElement): HTMLElement {
    const spinner_container = NewToDoFormContainer.querySelector('.spinner-conatiner') as HTMLElement
    const NewToDoForm = NewToDoFormContainer.querySelector('.add-new-todo') as HTMLElement
    (NewToDoForm).addEventListener('submit', async (e) => {
        e.preventDefault();
        const newTitle = (NewToDoForm.querySelector('input[name="title"]') as HTMLInputElement).value
        const newCategory = (NewToDoForm.querySelector('input[name="category"]') as HTMLInputElement).value
        const newDescription = (NewToDoForm.querySelector('textarea[name="description"]') as HTMLInputElement).value
        const newToDo: IToDo = {
            title: newTitle,
            category: newCategory,
            description: newDescription
        }

        if (!newTitle || !newCategory || !newDescription) {
            NewToDoFormContainer.querySelector('.warning-txt')?.classList.toggle('show')
            return
        }

        try {
            spinner_container.classList.toggle('show')
            await addToDoFunc(newToDo);
            (NewToDoForm as HTMLFormElement).reset()
            NewToDoFormContainer.classList.toggle('show')
        }
        catch (err) {
            console.log(err)
        }
        finally {
            spinner_container.classList.toggle('show')
        }
    })
    return NewToDoFormContainer
}

function addNewTodoFormContainerEvents(NewToDoFormContainer: HTMLElement): HTMLElement {
    let NewToDoForm = (NewToDoFormContainer.querySelector('.add-new-todo') as HTMLElement)
    NewToDoFormContainer = addCreateNewToDoFormContainerEvent(NewToDoFormContainer);
    NewToDoForm = addRemoveWarningTxtEvent(NewToDoForm)
    return NewToDoFormContainer
}



//Add the change state event
function addChangeStateEvent() {
    const todo_states = document.querySelector('.todo-states') as HTMLElement
    (todo_states.querySelector('.all-state') as HTMLElement).addEventListener('click', (e) => {
        stateFilter = 'All'
        todo_states.className = 'todo-states all'
        renderToDo()
    });
    (todo_states.querySelector('.done-state') as HTMLElement).addEventListener('click', (e) => {
        stateFilter = 'Done'
        todo_states.className = 'todo-states done'
        renderToDo()
    });
    (todo_states.querySelector('.undone-state') as HTMLElement).addEventListener('click', (e) => {
        stateFilter = 'Undone'
        todo_states.className = 'todo-states undone'
        renderToDo()
    })
}

//Add search by title and category event
function addSearchByTitleAndCategoryEvent() {
    (document.querySelector('.search-section') as HTMLElement).addEventListener('submit', async (e) => {
        e.preventDefault();
        const title = (document.querySelector('#search-title') as HTMLInputElement).value
        const category = (document.querySelector('#search-category') as HTMLInputElement).value
        try {
            (document.querySelector('.todo-section>.spinner-conatiner') as HTMLElement).classList.toggle('show')
            await searchToDoFunc(title, category)
        }
        catch (err) {
            console.log(err)
        }
        finally {
            (document.querySelector('.todo-section>.spinner-conatiner') as HTMLElement).classList.toggle('show')
        }
    })
}

//Add initial load
async function initialLoad() {
    try {
        (document.querySelector('.todo-section>.spinner-conatiner') as HTMLElement).classList.toggle('show')
        await loadToDoFunc()
    }
    catch (err) {
        console.log(err)
    }
    finally {
        (document.querySelector('.todo-section>.spinner-conatiner') as HTMLElement).classList.toggle('show')
    }
}

//Add the create todo formModel functionality
function addExitEvents() {
    const add_new_todo_container = document.querySelector(".add-new-todo-container") as HTMLElement
    const exit_icon = document.querySelector(".add-new-todo-container .exit-icon") as HTMLElement

    exit_icon.addEventListener('click', () => {
        add_new_todo_container.classList.toggle("show")
    })
}
addNewTodoFormContainerEvents(document.querySelector('.add-new-todo-container') as HTMLElement)
addExitEvents()
addChangeStateEvent()
addSearchByTitleAndCategoryEvent()
initialLoad()
