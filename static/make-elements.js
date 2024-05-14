// Returns a DOM element, spcified by inputs as follows
    // type: div, button
    // name: id
    // classNames: class names
    // contents: innerHTML
    // link: link
export function makeElement(type, name, classNames, contents, link) {
    const element = document.createElement(type)
    if(name) {
        element.id = name
    }
    if(classNames) {
        if(Array.isArray(classNames)) {
            element.classList.add(...classNames)
        } else {
            element.classList.add(classNames)
        }
    }
    if(contents) {
        element.innerHTML = contents
    }
    if(link) {
        element.href = link
    }

    return element
}