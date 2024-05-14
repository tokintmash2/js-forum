import { makeElement } from "./make-elements.js"
import { loginData } from "./fetch.js"

export function navbar() {

    const navSpan = document.getElementsByTagName('nav')[0]

    navSpan.innerHTML = ''

    const home = makeElement('a', '', '', 'Home', '/')
    navSpan.appendChild(home)

    loginData()
        .then(userData => {
            if(userData.LoggedIn) {
                const welcome =  makeElement('span', '', 'welcome', `Welcome, ${userData.username}!`, '')
                const createPost = makeElement('a', '', '', 'Create a new post', '/create-post')
                const likedPosts = makeElement('a', '', '', 'My liked posts', '/liked-posts')
                const myPosts = makeElement('a', '', '', 'My posts', '/my-posts')
                const logOut = makeElement('a', '', '', 'Log out', '/logout')
                navSpan.appendChild(welcome)
                navSpan.appendChild(createPost)
                navSpan.appendChild(likedPosts)
                navSpan.appendChild(myPosts)
                navSpan.appendChild(logOut)
            } else {
                const signUp = makeElement('a', '', '', 'Sign up', '/sign-up')
                const signIn = makeElement('a', '', '', 'Sign in', '/sign-in')
                navSpan.appendChild(signUp)
                navSpan.appendChild(signIn)
            }
        })
        .catch(error => {
            console.error('Error reading login data in nevbar.js', error)
        })
    
}