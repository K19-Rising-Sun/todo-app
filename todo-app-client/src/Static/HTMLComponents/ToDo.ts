export default function HTMLTodo(title: string,category: string,description: string,isCompleted: boolean): HTMLElement{
    const HTMLTodo = document.createElement('div')
    HTMLTodo.className=`todo ${isCompleted?'complete':''}`
    HTMLTodo.innerHTML=`
        <div class="spinner-conatiner">
            <div class="spinner"></div>
        </div>
        <p class="category">${category}</p>
        <div class="todo-header">
            <div class="todo-title"><h2>${title}</h2></div>
            <span class="icons">
                <span class="material-symbols-outlined edit-icon">
                    edit
                </span>
                <span class="material-symbols-outlined done-icon">
                    done
                    </span>
                <span class="material-symbols-outlined delete-icon">
                    delete
                </span>
            </span>
        </div>
        <p class="description">${description}</p>
    `

    return HTMLTodo
}
