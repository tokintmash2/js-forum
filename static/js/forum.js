import { fetchCatPosts, fetchOnlineUsers } from './fetch.js';
import { handleLogout } from "./JS-handlers.js"
import { makeElements } from './make-elements.js';
import { renderCreatePost, renderHomePage, signIn, signUp } from './render.js';
import { navbar } from './navbar.js';


// Get appDiv
let appDiv = document.getElementById('app');
const mainContent = makeElements({type: 'div', classNames: 'main-content'})
// const chatSection = makeElements({type: 'div', classNames: 'sidebar'})
appDiv.innerHTML = '';
// mainContent.innerHTML = '';

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
                fetchOnlineUsers()
                // connectWebSocket()
                break;
            case `/category/${categoryId}`:
            case `/like`:
                renderCategoryPage(categoryId)
                fetchOnlineUsers()

                break;
            case '/sign-in':
                signIn();
                break;
            case '/sign-up':
                signUp()
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

    // renderCreatePost()

    function renderCategoryPage(id) {
        appDiv.innerHTML = ''
        const pageTitle = makeElements({ type: 'h2', name: 'Title', contents: 'Posts in this Category' })
        mainContent.appendChild(pageTitle)
        fetchCatPosts(mainContent, id)
        appDiv.appendChild(mainContent)
    }

    // signIn()

    document.addEventListener('click', function (event) {
        if (event.target.tagName === 'A') {
            event.preventDefault();
            const path = event.target.getAttribute('href');
            window.history.pushState({}, '', path);
            handleRoute();
        }
    });

    // Handle browser navigation events
    window.onpopstate = function (event) {
        console.log('Popstate triggered:', window.location.pathname);
        handleRoute(window.location.pathname);
    };

    // Handle the initial route
    handleRoute();
});

