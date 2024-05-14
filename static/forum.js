import { fetchCategories, fetchRecentPosts } from './fetch.js';
import { makeElement } from './make-elements.js';
import { navbar } from './navbar.js';

// Get appDiv
let appDiv = document.getElementById('app');
appDiv.innerHTML = '';

// Main function
document.addEventListener('DOMContentLoaded', function () {

// renderHomePage()

function handleRoute(path) {
    console.log('From handleRoute:', window.location.pathname);
    switch (path) {
        case '/':
        case '':
            renderHomePage();
            break;
        case '/category/1':
            renderCategoryPage()
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

function renderCategoryPage() {

    console.log('Rendering Category Page');
    console.log('from renderCategoryPage')
    appDiv.innerHTML = '<p>This is the category page.</p>'; 

   
    const testDiv = document.createElement('div')
    testDiv.classList.add('testing')
    appDiv.appendChild(testDiv)
}

document.addEventListener('click', function(event) {
    console.log(event)
    if (event.target.tagName === 'A') {
        event.preventDefault();
        const path = event.target.getAttribute('href');
        console.log('Link clicked:', path);
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
handleRoute(window.location.pathname);

  
});




