import { fetchCatPosts, loginData } from './fetch.js'
import { handleLogout, websocket } from './JS-handlers.js'
import { $, makeElements } from './utils.js'
import { CreatePost } from './pages/CreatePost.js'
import { LogIn } from './pages/LogIn.js'
import { Register } from './pages/Register.js'
import { navbar } from './navbar.js'
import { Chat } from './components/Chat.js'
import { Sidebar } from './components/Sidebar.js'
import { Home } from './pages/Home.js'

export const App = {
	element: $('#app'),
	UUID: null,
	user: {
		username: null,
		id: null,
		loggedIn: false,
	},
	register_page: null,
	logIn_page: null,
	home_page: null,
	createPost_page: null,

	chatInstances: [],
	initSidebar: function () {
		this.sidebar = new Sidebar()
	},

	renderLogInPage: function () {
		if (!this.logIn_page) {
			this.logIn_page = new LogIn()
			this.logIn_page.render()
		} else {
			this.logIn_page.render()
		}
	},
	renderHomePage: function () {
		if (!this.home_page) {
			this.home_page = new Home()
			this.home_page.render()
		} else {
			this.home_page.render()
		}
	},
	renderRegisterPage: function () {
		if (!this.register_page) {
			this.register_page = new Register()
			this.register_page.render()
		} else {
			this.register_page.render()
		}
	},
	renderCreatePostPage: function () {
		if (!this.createPost_page) {
			this.createPost_page = new CreatePost()
			this.createPost_page.render()
		} else {
			this.createPost_page.render()
		}
	},

	setUp: async function () {
		this.UUID = localStorage.getItem('UUID')
		try {
			const userData = await loginData()
			console.log('setUp: userData', userData)
			this.user = {
				id: userData.userId,
				username: userData.username,
				loggedIn: userData.LoggedIn,
			}
			this.initSidebar()

			if (this.UUID) {
				if (!this.websocketConnected) {
					websocket()
				} else {
				}
			}
			return userData
		} catch (error) {
			console.error('Error loading the loginData:', error)
			throw error // Re-throw the error to maintain the promise chain
		}
	},
	initChat: function (partnetId, partnerUsername) {
		// check if the chat between you and partner already exists
		let chatInstance = this.chatInstances.find((chat) => chat.partnerId == partnetId)
		console.log('Init: chatInstance', chatInstance)
		if (!chatInstance) {
			console.log('!chatInstance')
			chatInstance = new Chat(partnetId, partnerUsername)
			this.chatInstances.push(chatInstance)
			chatInstance.render()
			chatInstance.show()
			chatInstance.loadConversation()
		}

		chatInstance.show()
	},

	renderCategoryPage: function (id) {
		console.log('renderCategoryPage', id)
		let mainContent = $('.main-content')
		if (!mainContent) {
			mainContent = makeElements({
				type: 'div',
				classNames: 'main-content',
			})
			this.element.appendChild(mainContent)
		} else {
			mainContent.innerHTML = ''
		}
		fetchCatPosts(this.element, id)

		navbar()
		this.sidebar.render()
	},

	renderNotFoundPage: function () {
		this.element.innerHTML = ''
		navbar()
		const mainContent = makeElements({ type: 'div', classNames: 'main-content' })
		this.element.appendChild(mainContent)
		const notFound = makeElements({
			type: 'h1',
			contents: 'Page not found',
		})
		mainContent.appendChild(notFound)
	},

	handleRoute: function () {
		console.log('handleRoute')
		const path = window.location.pathname
		const pathSegments = path.split('/').filter((segment) => segment.length > 0)
		const categoryId = parseInt(pathSegments[1], 10)

		App.setUp().then(() => {
			switch (path) {
				case '/':
				case '':
					App.user.loggedIn ? this.renderHomePage() : this.renderLogInPage()
					break
				case `/category/${categoryId}`:
				case `/like`:
					App.user.loggedIn ? this.renderCategoryPage(categoryId) : this.renderLogInPage()
					break
				case '/log-in':
					this.renderLogInPage()
					break
				case '/register':
					this.renderRegisterPage()
					break
				case '/api/logout':
					handleLogout()
					break
				case '/createPost':
					App.user.loggedIn ? this.renderCreatePostPage() : this.renderLogInPage()
					break
				case '/my-posts':
					App.user.loggedIn ? renderMyPosts() : this.renderLogInPage()
					break
				case '/liked-posts':
					App.user.loggedIn ? renderLikedPosts() : this.renderLogInPage()
					break
				default:
					this.renderNotFoundPage()
			}
		})
	},
}

// Main function
document.addEventListener('DOMContentLoaded', function () {
	App.handleRoute()

	document.addEventListener('click', function (event) {
		if (event.target.tagName === 'A') {
			event.preventDefault()
			const path = event.target.getAttribute('href')
			console.log('pushState triggered. Path:', path)
			window.history.pushState({}, '', path)
			App.handleRoute()
		}
	})

	// Handle browser navigation events
	window.onpopstate = function (event) {
		console.log('Popstate triggered:', window.location.pathname)
		App.handleRoute(window.location.pathname)
	}
})
