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

/**
 * @typedef {Object} ElementData
 * @property {string} type - The type of the element (e.g., 'div', 'a').
 * @property {string} [name] - The id of the element.
 * @property {string|string[]} [classNames] - The class or classes of the element.
 * @property {string} [contents] - The innerHTML of the element.
 * @property {string} [link] - The href of the element if it is a link.
 */

/**
 * 
 * @param {ElementData} data 
 * @returns {HTMLElement}
 */

export function makeElements(data) {
    const element = document.createElement(data.type)
    if(data.name) {
        element.id = data.name
    }
    if(data.classNames) {
        if(Array.isArray(data.classNames)) {
            element.classList.add(...data.classNames)
        } else {
            element.classList.add(data.classNames)
        }
    }
    if(data.contents) {
        element.innerHTML = data.contents
    }
    if(data.link) {
        element.href = data.link
    }

    return element
}