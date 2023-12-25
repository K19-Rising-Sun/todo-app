export default function HTMLEditTodo(title: string | null,category: string | null,description: string | null):HTMLElement{
    const HTMLEditTodo = document.createElement('form')
    HTMLEditTodo.className=`todo`
    HTMLEditTodo.innerHTML=`
        <div class="spinner-conatiner">
            <div class="spinner"></div>
        </div>
        <input type="text" name="category" placeholder="Category">

        <div class="todo-header">
            <input type="text" name="title" placeholder="Title">
            <input type="submit" value="Change" class="submit-btn">
        </div>
        <textarea name="description" cols="30" rows="10">${description}</textarea>
        <p class="warning-txt">Please fill in all field</p>
    `;
    ((HTMLEditTodo.querySelector('input[name="title"]') as HTMLInputElement).value as string|null)=title;
    ((HTMLEditTodo.querySelector('input[name="category"]') as HTMLInputElement).value as string|null)=category
    return HTMLEditTodo
}
