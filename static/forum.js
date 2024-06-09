import { fetchCategories, fetchRecentPosts, fetchCatPosts, makeCatChkboxes, fetchMyPosts, fetchLikedPosts } from './fetch.js';
import { handleLikeDislike, handleLogout, creatPostHandler } from "./JS-handlers.js"
import { makeElement } from './make-elements.js';
import { renderHomePage } from './render.js';
import { navbar } from './navbar.js';


// Get appDiv
let appDiv = document.getElementById('app');
appDiv.innerHTML = '';

// Main function
document.addEventListener('DOMContentLoaded', function () {

    

// renderHomePage()

function handleRoute() {
   
    const path = window.location.pathname
    const pathSegments = path.split('/').filter(segment => segment.length > 0)
    const categoryId = parseInt(pathSegments[1], 10)

    switch (path) {
        case '/':
        case '':
            renderHomePage();
            break;
        case `/category/${categoryId}`:
        case `/like`:
            renderCategoryPage(categoryId)
            break;
        case '/sign-in':
            signIn();
            break;
        case '/api/logout':
            handleLogout()
            renderHomePage()
            break
        case '/createPost':
            renderCreatePost()
            break
        case '/my-posts':
            renderMyPosts()
            break
        case '/liked-posts':
            renderLikedPosts()
            break
        default:
            break;
    }
}

function renderMyPosts() {
    appDiv.innerHTML = '';
    navbar()
    fetchMyPosts(appDiv) 
}

function renderLikedPosts() {
    appDiv.innerHTML = '';
    navbar()
    fetchLikedPosts(appDiv) 
}

// Render new Post page
function renderCreatePost() {
    appDiv.innerHTML = '';
    navbar()
    const newPostDiv = makeElement('div', 'newPostDiv', 'container', '', '')
    const pageTitle = makeElement('h1', '', '', 'Create a New Post', '')
    const newPostTitle = makeElement('div', 'newPostTitle', '', '', '', '')
    newPostTitle.innerHTML = `
        <label for="title">Title:</label><br>
        <input type="text" id="title" name="title"><br><br>
        `
    const newPostBody = makeElement('div', 'newPostBody', '', '', '', '')
    newPostBody.innerHTML = `
        <label for="content">Content:</label><br>
        <textarea id="content" name="content" rows="14" cols="50"></textarea><br><br>
        `
    const categorySelect = makeElement('div', 'catSelect', '', '', '', '')
    makeCatChkboxes(categorySelect)

    const submitBtn = makeElement('button', '', 'subm-button', 'Create Post', '')
    submitBtn.addEventListener('click', creatPostHandler);
    
    newPostDiv.appendChild(pageTitle)
    newPostDiv.appendChild(newPostTitle)
    newPostDiv.appendChild(newPostBody)
    newPostDiv.appendChild(categorySelect)
    newPostDiv.appendChild(submitBtn)
    appDiv.appendChild(newPostDiv)
}

function renderCategoryPage(id) {
    appDiv.innerHTML = '';
    const pageTitle = makeElement('h2', 'Title', '', 'Posts in this Category', '')
    appDiv.appendChild(pageTitle)
    fetchCatPosts(appDiv, id)
}

function signIn() {
    appDiv.innerHTML = `
    <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sign In</title>
    <link rel="stylesheet" href="/static/sign_in.css">
    <link rel="stylesheet" href="/static/GG_buttons.css">
    <form action="/sign-in">
        <div class="container">
            <h1>Sign In</h1>
            <p>{{.}}</p>
            <p>Please fill in this form to sign in.</p>
            <hr>

            <label><b>Email</b></label>
            <input type="text" placeholder="Enter Email" name="email" required>

            <label><b>Password</b></label>
            <input type="password" placeholder="Enter Password" name="password" required>
            
            <div class="clearfix">
                <button type="submit" class="signinbtn">Sign In</button>
            </div>
        </div>
    </form>
</head>
<body>
    <div class="oauth-buttons">
        <a href="/google-login" class="oauth-button google">
            <img src="/static/img/google.png" alt="Google Icon"> Sign in with Google
        </a>
        <a href="/github-login" class="oauth-button github">
            <img src="/static/img/github.png" alt="GitHub Icon"> Sign in with GitHub
        </a>
    </div>
</body>`
}

document.addEventListener('click', function(event) {
    if (event.target.tagName === 'A') {
        event.preventDefault();
        const path = event.target.getAttribute('href');
        window.history.pushState({}, '', path);
        handleRoute(path);
    }
});

// Handle browser navigation events
window.onpopstate = function(event) {
    console.log('Popstate triggered:', window.location.pathname);
    handleRoute(window.location.pathname);
};

// Handle the initial route
handleRoute();  
});

