import { makeElements } from "./make-elements.js"
import { loginData } from "./fetch.js"

export function navbar() {

    const navSpan = document.getElementsByTagName('nav')[0]

    navSpan.innerHTML = ''

    const home = makeElements({type:'a', contents: 'Home', link: '/'})
    navSpan.appendChild(home)

    loginData()
        .then(userData => {
            if(userData.LoggedIn) {
                const welcome =  makeElements({type:'span', classNames:'welcome', contents: `Welcome, ${userData.username}!`})
                const createPost = makeElements({type:'a', contents: 'Create a new post', link: '/createPost'})
                const logOut = makeElements({type: 'a', contents: 'Log out', link: '/api/logout'})
                navSpan.appendChild(welcome)
                navSpan.appendChild(createPost)
                navSpan.appendChild(logOut)
            } else {
                const signUp = makeElements({type:'a', contents: 'Sign up', link: '/sign-up'})
                const signIn = makeElements({type: 'a', contents: 'Sign in', link: '/sign-in'})
                navSpan.appendChild(signUp)
                navSpan.appendChild(signIn)
            }
        })
        .catch(error => {
            console.error('Error reading login data in nevbar.js', error)
        })
    
}