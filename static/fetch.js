import { makeElement } from "./make-elements.js";

export function fetchCategories(list) {
    fetch('/api/categories')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(categories => {
            categories.forEach(category => {
                const a = document.createElement('a');
                a.textContent = category.Name;
                a.href = `/category/${category.ID}`;
                list.appendChild(a);
            });
        })
        .catch(error => {
            console.error('Error loading the categories:', error);
            list.textContent = 'Failed to load categories.'; // Update text content directly on 'list'
        });
}

export function fetchCatPosts(body, id) {
    fetch(`/api/category/${id}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            const posts = data.CategoryPosts;
            if (Array.isArray(posts) && posts.length > 0) {
                posts.forEach(post => {
                    // Posts
                    const postDiv = makeElement('div', 'JSPost', 'post', '', '');
    
                    const postBox = makeElement('div', '', 'content-box', `
                        <h3>${post.Title}</h3><p>${post.Content}</p>`,
                        '');
                    postDiv.appendChild(postBox);
    
                    const postDetails = makeElement('div', '', 'post-details', `
                        <p>Author: ${post.Author}</p>
                        <p>Likes: ${post.Likes}, Dislikes: ${post.Dislikes}</p>
                        <p>Created at: ${new Date(post.CreatedAt).toLocaleDateString()}</p>`,
                        '');
                    postBox.appendChild(postDetails);

                    // Like buttons
                    if(data.LoggedIn) {
                        const likeButtons = makeElement('div', '', 'post-buttons', '', '')
                        likeButtons.innerHTML = `
                        <form action="/like" method="post">
                            <input type="hidden" name="post_id" value=${post.ID}>
                            <button type="submit" name="action" value="like">Like</button>
                        </form>
                        <form action="/dislike" method="post">
                            <input type="hidden" name="post_id" value=${post.ID}>
                            <button type="submit" name="action" value="dislike">Dislike</button>
                        </form>
                        `
                        postBox.appendChild(likeButtons)
                    }
    
                    // Comments
                    const commentBox = makeElement('div', '', 'comment-container', '', '');
                    const comments = makeElement('div', '', 'comments', '', '');
    
                    if (Array.isArray(post.Comments)) {
                        post.Comments.forEach(comment => {
                            const commentDiv = makeElement('div', '', 'comment', `
                                <p>${comment.Content}</p>
                                <p>Author: ${comment.Author}</p>
                                <p>Likes: ${comment.Likes}, Dislikes: ${comment.Dislikes}</p>`,
                                '');
                            comments.appendChild(commentDiv); // This goes inside ForEach
                        });
                    }
                    commentBox.appendChild(comments);
                    postBox.appendChild(commentBox);
                    body.appendChild(postDiv);
                });
            } else {
                const noPosts = makeElement('div', '', '', 'No posts in this category.', '')
                body.appendChild(noPosts);
            }
        })
        .catch(error => {
            console.error('Error loading recent posts:', error);
            body.textContent = 'Failed to load recent posts.'; // Update text content directly on 'body'
        });
}



export function fetchRecentPosts(body) {
    fetch('/api/recents')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(posts => {
            posts.forEach(post => {
                // Posts
                const postDiv = makeElement('div', 'JSPost', 'post', '', '')

                const postBox = makeElement('div', '', 'content-box', `
                    <h3>${post.Title}</h3><p>${post.Content}</p>`,
                    '')
                postDiv.appendChild(postBox);

                const postDetails = makeElement('div', '', 'post-details', `
                    <p>Author: ${post.Author}</p>
                    <p>Likes: ${post.Likes}, Dislikes: ${post.Dislikes}</p>
                    <p>Created at: ${new Date(post.CreatedAt).toLocaleDateString()}</p>`,
                    '')
                postBox.appendChild(postDetails)

                // Comments
                const commentBox = makeElement('div', '', 'comment-container', '', '')
                const comments = makeElement('div', '', 'comments', '', '')

                if (Array.isArray(post.Comments)) {
                    post.Comments.forEach(comment => {
                        const commentDiv = makeElement('div', '', 'comment', `
                        <p>${comment.Content}</p>
                        <p>Author: ${comment.Author}</p>
                        <p>Likes: ${comment.Likes}, Dislikes: ${comment.Dislikes}</p>`,
                            '')
                        comments.appendChild(commentDiv) // This goes inside ForEach
                    })
                }
                commentBox.appendChild(comments)
                postBox.appendChild(commentBox)
                body.appendChild(postDiv)
            })
        })
        .catch(error => {
            console.error('Error loadinng recent posts:', error);
            list.textContent = 'Failed to load recent posts.'; // Update text content directly on 'list'
        });
}

export function loginData() {
    return new Promise((resolve, reject) => {
        fetch('/api/home')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {

                let userData = {
                    LoggedIn: data.LoggedIn,
                    username: data.Username
                }
                resolve(userData)
            })
            .catch(error => {
                console.error('Error loading the loginData:', error);
                throw error
            });
    })
}