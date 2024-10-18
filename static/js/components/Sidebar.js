import { makeElements, $ } from '../utils.js'
import { App } from '../index.js'

export class Sidebar {
	render() {
		if (!$('#sidebar')) {
			const sidebar = makeElements({
				type: 'div',
				name: 'sidebar',
				classNames: 'sidebar',
				children: [
					makeElements({ type: 'h3', contents: 'Online contacts' }),
					makeElements({ type: 'div', classNames: 'online-users' }),
				],
			})
			App.element.appendChild(sidebar)
		}
	}
	appendOnlineUsers(users) {
		const onlineUsersElement = $('#sidebar > .online-users')

		if (!$('#sidebar')) {
			this.render()
		} else {
			onlineUsersElement.innerHTML = ''
		}

		if (users) {
			users.forEach((user) => {
				if (user.ID !== App.user.id) {
					const onlineUser = makeElements({
						type: 'div',
						classNames: 'user',
						contents: `${user.Username}`,
					})
					onlineUser.addEventListener('click', () => App.initChat(user.ID, user.Username))
					onlineUsersElement.appendChild(onlineUser)
				}
			})
		} else {
			const onlineUser = makeElements({
				type: 'div',
				contents: `No online users`,
			})
			onlineUsersElement.appendChild(onlineUser)
		}
	}
}
