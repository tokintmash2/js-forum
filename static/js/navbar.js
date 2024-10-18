import { makeElements, $ } from './utils.js'
import { App } from './index.js'

export function navbar() {
	const navSpan = $('nav')

	navSpan.innerHTML = ''

	const brand = makeElements({ type: 'a', classNames: 'brand', link: '/' })
	navSpan.appendChild(brand)

	let actions = $('header .actions')
	if (!actions) {
		actions = makeElements({
			type: 'div',
			classNames: ['actions', 'flex'],
		})
		navSpan.appendChild(actions)
	}

	actions.innerHTML = ''

	if (App.user.loggedIn) {
		console.log('logged in')
		const createNewPost = makeElements({
			type: 'a',
			contents: 'Create a new post',
			link: '/createPost',
		})

		const logOut = makeElements({ type: 'a', contents: 'Log out', link: '/api/logout' })

		navSpan.appendChild(createNewPost)
		navSpan.appendChild(logOut)
	} else {
		console.log('not logged in')
		// const logIn = makeElements({ type: 'a', contents: 'Log in', link: '/log-in' })
		const logIn = makeElements({ type: 'a', contents: 'Log in', link: '/sign-in' })
		const signUp = makeElements({ type: 'a', contents: 'Register', link: '/register' })
		navSpan.appendChild(logIn)
		navSpan.appendChild(signUp)
	}
}
