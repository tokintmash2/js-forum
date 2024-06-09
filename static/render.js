import { fetchCategories, fetchRecentPosts } from "./fetch.js";
import { makeElement } from "./make-elements.js";
import { navbar } from "./navbar.js";

// Get appDiv
let appDiv = document.getElementById('app');
appDiv.innerHTML = '';

export function renderHomePage() {
    appDiv.innerHTML = '';
    navbar()

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