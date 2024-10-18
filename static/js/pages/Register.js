import { navbar } from '../navbar.js'
import { makeElements, $, validateEmail, clearPage, clearMainContentContainer } from '../utils.js'
import { App } from '../index.js'

export class Register {
	constructor() {}
	render() {
		clearPage()
		navbar()
		clearMainContentContainer()
		const formContainer = makeElements({ type: 'div', classNames: 'container' })
		const pageTitle = makeElements({ type: 'h1', contents: 'Register new user' })
		const usernameField = makeElements({
			type: 'div',
			name: 'username',
			children: [
				makeElements({ type: 'label', contents: '<b>Username</b>' }),
				makeElements({
					type: 'input',
					attributes: {
						type: 'text',
						placeholder: 'Enter Username',
						id: 'signup-username',
						required: 'true',
					},
				}),
			],
		})
		const emailField = makeElements({
			type: 'div',
			name: 'email',
			children: [
				makeElements({ type: 'label', contents: '<b>Email</b>' }),
				makeElements({
					type: 'input',
					attributes: {
						type: 'email',
						placeholder: 'Enter email',
						id: 'signup-email',
						required: 'true',
					},
				}),
			],
		})

		const firstNameField = makeElements({
			type: 'div',
			name: 'firstName',
			children: [
				makeElements({ type: 'label', contents: '<b>First Name</b>' }),
				makeElements({
					type: 'input',
					attributes: {
						type: 'text',
						placeholder: 'Enter first name',
						id: 'signup-firstName',
						required: 'true',
					},
				}),
			],
		})

		const lastNameField = makeElements({
			type: 'div',
			name: 'lastName',
			children: [
				makeElements({ type: 'label', contents: '<b>Last Name</b>' }),
				makeElements({
					type: 'input',
					attributes: {
						type: 'text',
						placeholder: 'Enter last name',
						id: 'signup-lastName',
						required: 'true',
					},
				}),
			],
		})

		const ageField = makeElements({
			type: 'div',
			name: 'age',
			children: [
				makeElements({ type: 'label', contents: '<b>Age</b>' }),
				makeElements({
					type: 'input',
					attributes: {
						type: 'number',
						placeholder: 'Enter your age',
						id: 'signup-age',
						required: 'true',
					},
				}),
			],
		})

		const genderField = makeElements({
			type: 'div',
			name: 'gender',
			children: [
				makeElements({
					type: 'label',
					contents: '<b>Gender</b>',
					attributes: { for: 'signup-gender' },
				}),
				makeElements({
					type: 'select',
					attributes: {
						id: 'signup-gender',
						name: 'signup-gender',
					},
					children: [
						makeElements({
							type: 'option',
							contents: 'Male',
							attributes: { selected: 'true', value: 'male' },
						}),
						makeElements({
							type: 'option',
							contents: 'Female',
							attributes: { value: 'female' },
						}),
						makeElements({
							type: 'option',
							contents: 'Other',
							attributes: { value: 'other' },
						}),
					],
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
						id: 'signup-password',
						required: 'true',
					},
				}),
			],
		})

		const submitBtn = makeElements({
			type: 'button',
			classNames: 'subm-button',
			contents: 'Register',
		})
		submitBtn.addEventListener('click', () => {
			this.registerUser()
		})

		formContainer.appendChild(pageTitle)

		formContainer.appendChild(usernameField)
		formContainer.appendChild(firstNameField)
		formContainer.appendChild(lastNameField)
		formContainer.appendChild(ageField)
		formContainer.appendChild(genderField)
		formContainer.appendChild(emailField)
		formContainer.appendChild(passwordField)
		formContainer.appendChild(submitBtn)
		$('.main-content').appendChild(formContainer)
	}
	registerUser() {
		const username = $('#signup-username').value
		const firstName = $('#signup-firstName').value
		const lastName = $('#signup-lastName').value
		const age = parseInt($('#signup-age').value) // Convert age to an integer
		const gender = $('#signup-gender').value
		const email = $('#signup-email').value
		const password = $('#signup-password').value

		if (!username || !email || !password || !firstName || !lastName || !age || !gender) {
			alert('All fields are required.')
			return
		}

		if (!validateEmail(email)) {
			alert('Please use a valid email address!')
			return
		}

		fetch('/api/signup', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				username: username,
				email: email,
				password: password,
				firstName: firstName,
				lastName: lastName,
				age: age,
				gender: gender,
			}),
		})
			.then((response) => response.json())
			.then((data) => {
				if (data.success) {
					history.pushState({}, '', '/')
					App.handleRoute()
				} else {
					alert('Signup failed: ' + data.message)
				}
			})
			.catch((error) => console.error('Error:', error))
	}
}
