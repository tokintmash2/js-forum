import { App } from '../index.js'
import { makeCatChkboxes } from '../fetch.js'
import { makeElements, clearMainContentContainer, $ } from '../utils.js'
import { creatPostHandler } from '../JS-handlers.js'
import { navbar } from '../navbar.js'

export class CreatePost {
	render() {
		clearMainContentContainer()

		navbar()
		App.sidebar.render()

		const pageTitle = makeElements({ type: 'h1', contents: 'Create a New Post' })

		const newPostDiv = makeElements({
			type: 'div',
			name: 'newPostDiv',
			classNames: 'content-box',
		})

		const newPostTitle = makeElements({ type: 'div', name: 'newPostTitle' })
		newPostTitle.innerHTML = `
    <label for="title">Title:</label><br>
    <input type="text" id="title" name="title"><br><br>
    `
		const newPostBody = makeElements({ type: 'div', name: 'newPostBody' })
		newPostBody.innerHTML = `
    <label for="content">Content:</label><br>
    <textarea id="content" name="content" rows="14" cols="50"></textarea><br><br>
    `
		const categorySelect = makeElements({ type: 'div', name: 'catSelect' })
		makeCatChkboxes(categorySelect)

		const submitBtn = makeElements({
			type: 'button',
			classNames: 'subm-button',
			contents: 'Create Post',
		})
		submitBtn.addEventListener('click', creatPostHandler)

		$('.main-content').appendChild(pageTitle)
		newPostDiv.appendChild(newPostTitle)
		newPostDiv.appendChild(newPostBody)
		newPostDiv.appendChild(categorySelect)
		newPostDiv.appendChild(submitBtn)
		$('.main-content').appendChild(newPostDiv)
	}
}
