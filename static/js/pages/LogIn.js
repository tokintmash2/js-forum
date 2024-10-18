import { $, makeElements, validateEmail, clearPage, clearMainContentContainer } from '../utils.js'
import { websocket } from '../JS-handlers.js'
import { navbar } from '../navbar.js'
import { App } from '../index.js'

export class LogIn {
	constructor(UUID) {
		this.UUID = UUID
	}
	render() {
		clearPage()
		clearMainContentContainer()
		navbar()

		const formContainer = makeElements({ type: 'div', classNames: 'container' })
		const pageTitle = makeElements({ type: 'h1', contents: 'Log in' })
		const identifierField = makeElements({
			type: 'div',
			name: 'email',
			children: [
				makeElements({ type: 'label', contents: '<b>Email or username</b>' }),
				makeElements({
					type: 'input',
					attributes: {
						type: 'text',
						placeholder: 'Enter email or username',
						id: 'login-identifier',
						required: 'true',
					},
				}),
			],
		})
		const passwordField = makeElements({
			type: 'div',
			name: 'password',
			children: [
				makeElements({ type: 'label', contents: '<b>Password</b>' }),
				makeElements({
					type: 'input',
					attributes: {
						type: 'password',
						placeholder: 'Enter password',
						id: 'login-password',
						required: 'true',
					},
				}),
			],
		})

		const submitBtn = makeElements({
			type: 'button',
			classNames: 'subm-button',
			contents: 'Log in',
		})
		submitBtn.addEventListener('click', () => {
			this.authenticateUser()
		})

		formContainer.appendChild(pageTitle)
		formContainer.appendChild(identifierField)
		formContainer.appendChild(passwordField)
		formContainer.appendChild(submitBtn)
		$('.main-content').appendChild(formContainer)
	}

	async authenticateUser() {
		console.log('LogIn.authenticateUser')
		const identifier = $('#login-identifier').value
		const password = $('#login-password').value

		if (!identifier || !password) {
			alert('Both fields are required.')
			return
		}

		const loginResponse = await fetch('/api/login', {
			method: 'POST',
			body: JSON.stringify({ identifier: identifier, password: password }),
			headers: { 'Content-Type': 'application/json' },
		})

		const loginData = await loginResponse.json()
		if (loginData.success) {
			this.UUID = loginData.UUID // Assuming the server sends a UUID
			localStorage.setItem('UUID', this.UUID)
			console.log('UUID', this.UUID)
			console.log(loginData)

			// Wait for the cookie to be set
			setTimeout(() => {
				history.pushState({}, '', '/')
				App.handleRoute()
			}, 1000) // Adjust the timeout as necessary - set from 1000 to 2000 to see if it fixes the userID issue
		} else {
			alert('Login failed: ' + loginData.message)
		}
	}
}
