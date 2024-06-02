import { fetchCategories, fetchRecentPosts, fetchCatPosts } from './fetch.js';
import { handleLikeDislike } from "./JS-handlers.js"
import { makeElement } from './make-elements.js';
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
        default:
            break;
    }
}

function renderHomePage() {
    appDiv.innerHTML = '';
    navbar()

    const testLlink = makeElement('a', 'test', '', 'test', '/test')
    appDiv.appendChild(testLlink)

    // Create category headline
    const catHeader = makeElement('h2', 'JScatHeader', 'category-header', 'Categories')
    appDiv.appendChild(catHeader)
    // Create category list
    const list = makeElement('div', 'categoryList', 'categories-container', '')
    // Insert categories right after category header
    catHeader.insertAdjacentElement('afterend', list);
    // Craete recent posts headline
    const rcntPostsHeader = makeElement('h2', 'JSrcntsHeader', 'category-header', 'Recent posts')
    appDiv.appendChild(rcntPostsHeader)
    // Fetch categories
    fetchCategories(list)
    // Fetch recent posts
    fetchRecentPosts(appDiv)
}

function renderCategoryPage(id) {

    // User functionality missing
    // JS toggle script

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




