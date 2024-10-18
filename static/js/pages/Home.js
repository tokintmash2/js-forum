import { App } from '../index.js'
import { fetchCategories } from '../fetch.js'
import { makeElements, clearMainContentContainer, $ } from '../utils.js'
import { navbar } from '../navbar.js'

export class Home {
	constructor() {
		this.appElement = $('#app')
	}
	render() {
		clearMainContentContainer()
		console.log('about to render navbar. App.user:', App.user)
		navbar()
		App.sidebar.render()

		// Create category headline
		const catHeader = makeElements({
			type: 'h2',
			name: 'JScatHeader',
			classNames: 'category-header',
			contents: 'Categories',
		})
		$('.main-content').appendChild(catHeader)
		// Create category list
		const list = makeElements({
			type: 'div',
			name: 'categoryList',
			classNames: 'categories-container',
		})
		// Insert categories right after category header
		catHeader.insertAdjacentElement('afterend', list)

		const rcntPostsHeader = makeElements({
			type: 'h2',
			name: 'JSrcntsHeader',
			classNames: 'category-header',
			contents: 'Recent posts',
		})
		$('.main-content').appendChild(rcntPostsHeader)

		fetchCategories(list)
		this.fetchRecentPosts($('.main-content'))
	}
	fetchRecentPosts(body) {
		fetch('/api/recents')
			.then((response) => {
				if (!response.ok) {
					throw new Error('Network response was not ok')
				}
				return response.json()
			})
			.then((posts) => {
				posts.forEach((post) => {
					// Posts
					const postDiv = makeElements({
						type: 'div',
						name: `${post.ID}`,
						attributes: {
							'data-id': `post-${post.ID}`,
							class: `${
								!post.Comments || post.Comments.length === 0
									? 'no-comments post'
									: 'post'
							} `,
						},
					})

					const postBox = makeElements({
						type: 'div',
						classNames: 'content-box',
					})
					postDiv.appendChild(postBox)

					const postTextContainer = makeElements({
						type: 'div',
						contents: `
                        <h3>${post.Title}</h3><p>${post.Content}</p>`,
						attributes: {
							class: `${
								!post.Comments || post.Comments.length === 0
									? 'no-comments post-text-container'
									: 'post-text-container'
							} `,
						},
					})

					postBox.appendChild(postTextContainer)

					const postDetails = makeElements({
						type: 'div',
						classNames: 'post-details',
						contents: `
                            <p>Author: ${post.Author}</p>
                            <p>Likes: ${post.Likes}, Dislikes: ${post.Dislikes}</p>
                            <p>Created at: ${new Date(post.CreatedAt).toLocaleDateString()}</p>
                            <p class="amount-of-comments">
            ${post.Comments && post.Comments.length > 0 ? post.Comments.length : 0} comment(s)
        </p>`,
					})
					postTextContainer.appendChild(postDetails)

					// Comments
					const commentBox = makeElements({
						type: 'div',
						classNames: ['comment-container', 'hidden'],
					})
					const comments = makeElements({ type: 'div', classNames: 'comments' })

					if (Array.isArray(post.Comments)) {
						post.Comments.forEach((comment) => {
							const commentDiv = makeElements({
								type: 'div',
								classNames: 'comment',
								contents: `
                            <p>${comment.Content}</p>
                            <p>Author: ${comment.Author}</p>
                            <p>Likes: ${comment.Likes}, Dislikes: ${comment.Dislikes}</p>`,
							})
							comments.appendChild(commentDiv) // This goes inside ForEach
						})
					}

					commentBox.appendChild(comments)
					postBox.appendChild(commentBox)
					body.appendChild(postDiv)
					// Add event listener if post has comments
					if (post.Comments && post.Comments.length > 0) {
						const myPostElement = $(`.post[data-id=post-${post.ID}]`)
						const myPostTextContainerElement =
							myPostElement.querySelector('.post-text-container')
						myPostTextContainerElement.addEventListener('click', () => {
							myPostElement
								.querySelector('.comment-container')
								.classList.toggle('hidden')
						})
					}
				})
			})
			.catch((error) => {
				console.error('Error loadinng recent posts:', error)
				list.textContent = 'Failed to load recent posts.' // Update text content directly on 'list'
			})
	}
}
