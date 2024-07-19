import { fetchCategories, fetchRecentPosts, makeCatChkboxes } from "./fetch.js";
import { makeElements } from "./make-elements.js";
import { login, creatPostHandler, signup } from "./JS-handlers.js"
import { navbar } from "./navbar.js";

// Get appDiv
let appDiv = document.getElementById('app');
appDiv.innerHTML = '';

export function renderHomePage() {
    appDiv.innerHTML = '';
    navbar()
    const mainContent = makeElements({type: 'div', classNames: 'main-content'})
    const chatSection = makeElements({type: 'div', classNames: 'sidebar'})
    // Create category headline
    const catHeader = makeElements({type:'h2', name:'JScatHeader', classNames:'category-header', contents:'Categories'})
    mainContent.appendChild(catHeader)
    // Create category list
    const list = makeElements({type:'div', name:'categoryList', classNames:'categories-container'})
    // Insert categories right after category header
    catHeader.insertAdjacentElement('afterend', list);

    const rcntPostsHeader = makeElements({type:'h2', name:'JSrcntsHeader', classNames:'category-header', contents:'Recent posts'})
    mainContent.appendChild(rcntPostsHeader)    
    
    appDiv.appendChild(mainContent)
    appDiv.appendChild(chatSection)

    fetchCategories(list)
    fetchRecentPosts(mainContent)
}

export function signIn() {
    appDiv.innerHTML = '';
    const formContainer = makeElements({ type:'div', classNames:'container' })
    const pageTitle = makeElements({ type:'h1', contents:'Sing in' })
    const usernameField = makeElements({ type: 'div', name: 'email' })
    const passwordField = makeElements({ type: 'div', name: 'password' })

    usernameField.innerHTML = `
        <label><b>Email</b></label>
        <input type="email" placeholder="Enter Email" id="login-email" required>
            `
    passwordField.innerHTML = `
        <label><b>Password</b></label>
        <input type="password" placeholder="Enter Password" id="login-password" required>
            `
    const submitBtn = makeElements({ type: 'button', classNames: 'subm-button', contents: 'Sign in' })
    submitBtn.addEventListener('click', login);

    formContainer.appendChild(pageTitle)
    formContainer.appendChild(usernameField)
    formContainer.appendChild(passwordField)
    formContainer.appendChild(submitBtn)
    appDiv.appendChild(formContainer)
}

export function signUp() {
    appDiv.innerHTML = '';
    const formContainer = makeElements({ type:'div', classNames:'container' })
    const pageTitle = makeElements({ type:'h1', contents:'Sing up' })
    const usernameField = makeElements({ type: 'div', name: 'username' })
    const emailField = makeElements({ type: 'div', name: 'email' })
    const passwordField = makeElements({ type: 'div', name: 'password' })

    usernameField.innerHTML = `
        <label><b>Username</b></label>
        <input type="text" placeholder="Enter Username" id="signup-username" required>
        `

    emailField.innerHTML = `
        <label><b>Email</b></label>
        <input type="email" placeholder="Enter Email" id="signup-email" required>
            `
    passwordField.innerHTML = `
        <label><b>Password</b></label>
        <input type="password" placeholder="Enter Password" id="signup-password" required>
            `
    const submitBtn = makeElements({ type: 'button', classNames: 'subm-button', contents: 'Sign un' })
    submitBtn.addEventListener('click', signup);

    formContainer.appendChild(pageTitle)
    formContainer.appendChild(usernameField)
    formContainer.appendChild(emailField)
    formContainer.appendChild(passwordField)
    formContainer.appendChild(submitBtn)
    appDiv.appendChild(formContainer)
}

// Render new Post page
export function renderCreatePost() {
    appDiv.innerHTML = '';
    navbar()
    const newPostDiv = makeElements({ type: 'div', name: 'newPostDiv', classNames: 'container' })
    const pageTitle = makeElements({ type: 'h1', contents: 'Create a New Post' })
    const newPostTitle = makeElements({ type: 'div', name: 'newPostTitle' })
    newPostTitle.innerHTML = `
    <label for="title">Title:</label><br>
    <input type="text" id="title" name="title"><br><br>
    `
    const newPostBody = makeElements({ type: 'div', name: 'newPostBody' })
    newPostBody.innerHTML = `
    <label for="content">Content:</label><br>
    <textarea id="content" name="content" rows="14" cols="50"></textarea><br><br>
    `
    const categorySelect = makeElements({ type: 'div', name: 'catSelect' })
    makeCatChkboxes(categorySelect)

    const submitBtn = makeElements({ type: 'button', classNames: 'subm-button', contents: 'Create Post' })
    submitBtn.addEventListener('click', creatPostHandler);

    newPostDiv.appendChild(pageTitle)
    newPostDiv.appendChild(newPostTitle)
    newPostDiv.appendChild(newPostBody)
    newPostDiv.appendChild(categorySelect)
    newPostDiv.appendChild(submitBtn)
    appDiv.appendChild(newPostDiv)
}