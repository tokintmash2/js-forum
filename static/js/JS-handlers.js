import { App } from './index.js'
import { appendNewComment } from './append.js'
import { $ } from './utils.js'
import { Home } from './pages/Home.js'
import { Chat } from './components/Chat.js'

export let ws

export function websocket() {
	ws = new WebSocket('ws://localhost:4000/ws')
	ws.onopen = () => {
		console.log('WebSocket connection established')
		App.websocketConnected = true
	}

	ws.onmessage = function (event) {
		// console.log("Received event:", event)
		let message = JSON.parse(event.data)
		// console.log('Received message:', message)
		// console.log("Message keys:", Object.keys(message));
		if (message.type === 'online_users') {
			// console.log("Received message online_users:", message.content)
			App.sidebar.appendOnlineUsers(message.online_users)
		}
		if (message.type === 'conversation') {
			let chatInstance = App.chatInstances.find(
				(chat) => message.sender == chat.partnerId || message.recipient == chat.partnerId
			)
			console.log('allchatinstances', App.chatInstances)
			console.log('conversation', message)
			console.log('mychatinstance', chatInstance)
			chatInstance.appendConversation(message)
		}
		if (message.type === 'user-message') {
			console.log('Received "user-message":', message)

			// case 1: I am the sender -> chatInstance must exist and be visible -> look for chatInstantace with id chatArea-${message.recipient}
			// case 2: I am the recipient -> look for chatInstantace with id chatArea-${message.sender}
			// case 3: I am both -> chatInstance must exist and be visible
			let chatInstance = undefined
			if (message.sender == App.user.id) {
				// case 1 and/or case 3
				chatInstance = App.chatInstances.find((chat) => chat.partnerId == message.recipient)
			} else {
				// case 2 and/or case 3
				chatInstance = App.chatInstances.find((chat) => chat.partnerId == message.sender)
				if (!chatInstance) {
					chatInstance = new Chat(message.sender, message.sender_username)
					App.chatInstances.push(chatInstance)
					chatInstance.render()
					chatInstance.loadConversation()
				}
			}

			console.log('"user-message": mychatinstance', chatInstance)

			// Is it currently open?
			if (chatInstance.visible) {
				chatInstance.appendConversation(message)
			} else {
				chatInstance.show()
				chatInstance.appendConversation(message)
			}

			/* chatInstance.loadConversation()
			chatInstance.appendConversation(message) */
		}
	}

	ws.onerror = function (error) {
		console.error('WebSocket Error:', error)
	}

	ws.onclose = function (event) {
		ws = null
		App.websocketConnected = false
	}
}

function hasCookies() {
	console.log(document.cookie)
	return document.cookie.split(';').some((cookie) => cookie.trim().startsWith('session='))
}

export function handleLikeDislike(postId, isLike) {
	const url = isLike ? '/api/likePost' : '/api/dislikePost'
	console.log('Post ID:', postId)
	console.log('URL:', url)

	fetch(url, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({ post_id: postId }),
	})
		.then((response) => {
			console.log('Response Status:', response.status)

			if (!response.ok) {
				return response.text().then((text) => {
					throw new Error(text)
				})
			}
			return response.json()
		})

		.then((data) => {
			if (data.success) {
				console.log(data.newLikeCount) // New like count
				console.log(document.cookie)
				const postDetailsElement = document
					.querySelector(`.like-button[data-post-id="${postId}"]`)
					.closest('.post')
					.querySelector('.post-details')
				if (postDetailsElement) {
					const likesDislikesElement = Array.from(
						postDetailsElement.querySelectorAll('p')
					).find((p) => p.textContent.includes('Likes:'))
					if (likesDislikesElement) {
						// Extract dislikes count from the current text
						const dislikesCount =
							likesDislikesElement.textContent.match(/Dislikes: (\d+)/)[1]
						likesDislikesElement.textContent = `
                        Likes: ${data.newLikeCount}, Dislikes: ${data.newDislikeCount}`
					} else {
						console.error('Likes/Dislikes element not found.')
					}
				} else {
					console.error('Post details element not found.')
				}
			}
		})

		.catch((error) => {
			console.error('Error liking/disliking post:', error)
		})
}

export function handleCommentLikes(commentId, isLike) {
	const url = isLike ? '/api/likeComment' : '/api/dislikeComment'
	console.log('Comment ID:', commentId)
	console.log('URL:', url)

	fetch(url, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({ comment_id: commentId }),
	})
		.then((response) => {
			console.log('Response Status:', response.status)

			if (!response.ok) {
				return response.text().then((text) => {
					throw new Error(text)
				})
			}
			return response.json()
		})

		.then((data) => {
			if (data.success) {
				console.log(data.newLikeCount) // New like count
				const postDetailsElement = document
					.querySelector(`.like-button[data-comment-id="${commentId}"]`)
					.closest('.comment')
				if (postDetailsElement) {
					const likesDislikesElement = Array.from(
						postDetailsElement.querySelectorAll('p')
					).find((p) => p.textContent.includes('Likes:'))
					if (likesDislikesElement) {
						// Extract dislikes count from the current text
						const dislikesCount =
							likesDislikesElement.textContent.match(/Dislikes: (\d+)/)[1]
						likesDislikesElement.textContent = `
                        Likes: ${data.newLikeCount}, Dislikes: ${data.newDislikeCount}`
					} else {
						console.error('Likes/Dislikes element not found.')
					}
				} else {
					console.error('Post details element not found.')
				}
			}
		})

		.catch((error) => {
			console.error('Error liking/disliking post:', error)
		})
}

export function handleLogout() {
	fetch('/api/logout', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
	})
		.then((response) => response.json())
		.then((data) => {
			if (data.success) {
				//document.getElementById('app').innerHTML = ''
				App.UUID = undefined
				localStorage.removeItem('UUID')
				history.pushState({}, '', '/log-in')
				App.handleRoute()
			} else {
				console.error('Logout failed')
			}
		})
		.catch((error) => {
			console.error('Error during logout:', error)
		})
}

export function creatPostHandler() {
	const title = document.getElementById('title').value
	const content = document.getElementById('content').value
	const categoryIDs = Array.from(
		document.querySelectorAll('input[name="categoryIDs"]:checked')
	).map((checkbox) => parseInt(checkbox.value))

	const postData = {
		title: title,
		content: content,
		categoryIDs: categoryIDs,
	}

	fetch('/api/create-post', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(postData),
	})
		.then((response) => {
			if (!response.ok) {
				return response.text().then((text) => {
					throw new Error(text)
				})
			}
			return response.json()
		})
		.then((data) => {
			if (data.success) {
				history.pushState({}, '', '/')
				App.handleRoute()
				console.log('Post created successfully')
			} else {
				console.error('Error creating post')
			}
		})
		.catch((error) => {
			console.error('Error creating post:', error)
		})
}

export function createCommentHandler(post_ID) {
	const postID = parseInt(post_ID)
	const content = document.getElementById(`comment-post-id-${postID}`).value

	console.log(postID, content)

	const commentData = {
		post_id: postID,
		content: content,
	}

	console.log('Comment Data:', commentData)

	fetch('/api/add-comment', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(commentData),
	})
		.then((response) => {
			if (!response.ok) {
				return response.text().then((text) => {
					throw new Error(text)
				})
			}
			return response.json()
		})
		.then((data) => {
			if (data.success) {
				console.log('data:', data)
				console.log('Comment added successfully')
				appendNewComment(data.comment)
				document.getElementById(`comment-post-id-${postID}`).value = '' // Clear input
			} else {
				console.error('Error adding comment')
			}
		})
		.catch((error) => {
			console.error('Error adding comment:', error)
		})
}

export function toggleComments(postID) {
	console.log('toggleComments: id: ', postID)
	const commentsContainer = $('.post[name="1"]')
	//commentsContainer.style.display = commentsContainer.style.display === 'none' ? 'block' : 'none'
}
