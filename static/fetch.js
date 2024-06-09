import { makeElement } from "./make-elements.js";
import { handleLikeDislike, handleCommentLikes, createCommentHandler } from "./JS-handlers.js"
import { reinitializeCommentFunctionality } from "./commentToggle.js"

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

// Duplicates lots of the fetchCategories func. Would be nice to combine
export function makeCatChkboxes(container) {
    fetch('/api/categories')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(categories => {
            categories.forEach(category => {
                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkbox.id = `category_${category.ID}`;
                checkbox.name = 'categoryIDs';
                checkbox.value = category.ID;

                const label = document.createElement('label');
                label.htmlFor = `category_${category.ID}`;
                label.textContent = category.Name;

                const br = document.createElement('br');

                container.appendChild(checkbox);
                container.appendChild(label);
                container.appendChild(br);
            });
        })
        .catch(error => {
            console.error('Error loading the categories:', error);
            container.textContent = 'Failed to load categories.'; // Update text content directly on 'container'
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

                     if(data.LoggedIn) {

                        reinitializeCommentFunctionality()

                        const likeButtons = makeElement('div', '', 'post-buttons', '', '')
                        const addCommentDiv = makeElement('div', 'add-comment', 'add-comment-section', '', '')
                        
                        likeButtons.innerHTML = `
                            <button class="like-button" data-post-id="${post.ID}">Like</button>
                            <button class="dislike-button" data-post-id="${post.ID}">Dislike</button>
                            `
                        likeButtons.querySelector('.like-button').addEventListener('click', () => handleLikeDislike(post.ID, true));
                        likeButtons.querySelector('.dislike-button').addEventListener('click', () => handleLikeDislike(post.ID, false));

                        addCommentDiv.innerHTML = `
                            <button class="add-comment-btn">Write Comment</button>
                            <textarea name="comment_content" rows="6" cols="50" placeholder="Leave your comment here"></textarea>
                            <button class="sumbit-comment-btn" type="submit">Add Comment</button>
                            `
                        addCommentDiv.querySelector('.sumbit-comment-btn').addEventListener('click', () => {
                            console.log(post.ID)
                            createCommentHandler(post.ID)
                        } )

                        

                        postBox.appendChild(likeButtons);
                        postBox.appendChild(addCommentDiv)
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
                                if(data.LoggedIn) {
                                    const commentButtons = makeElement('div', '', 'comment-buttons', '', '')
                                    commentButtons.innerHTML = `
                                    <button class="like-button" data-comment-id="${comment.ID}">Like</button>
                                    <button class="dislike-button" data-comment-id="${comment.ID}">Dislike</button>
                                    `
                                    commentButtons.querySelector('.like-button').addEventListener('click', () => handleCommentLikes(comment.ID, true));
                                    commentButtons.querySelector('.dislike-button').addEventListener('click', () => handleCommentLikes(comment.ID, false));

                                    commentDiv.appendChild(commentButtons)
                                }
                            comments.appendChild(commentDiv);
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

export function fetchMyPosts(body) {
    fetch(`/api/my-posts`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            const posts = data.UserPosts;
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

                     if(data.LoggedIn) {
                        const likeButtons = makeElement('div', '', 'post-buttons', '', '')
                        
                        likeButtons.innerHTML = `
                            <button class="like-button" data-post-id="${post.ID}">Like</button>
                            <button class="dislike-button" data-post-id="${post.ID}">Dislike</button>
                            `
                        likeButtons.querySelector('.like-button').addEventListener('click', () => handleLikeDislike(post.ID, true));
                        likeButtons.querySelector('.dislike-button').addEventListener('click', () => handleLikeDislike(post.ID, false));
                        
                        postBox.appendChild(likeButtons);
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
                                if(data.LoggedIn) {
                                    const commentButtons = makeElement('div', '', 'comment-buttons', '', '')
                                    commentButtons.innerHTML = `
                                    <button class="like-button" data-comment-id="${comment.ID}">Like</button>
                                    <button class="dislike-button" data-comment-id="${comment.ID}">Dislike</button>
                                    `
                                    commentButtons.querySelector('.like-button').addEventListener('click', () => handleCommentLikes(comment.ID, true));
                                    commentButtons.querySelector('.dislike-button').addEventListener('click', () => handleCommentLikes(comment.ID, false));

                                    commentDiv.appendChild(commentButtons)
                                }
                            comments.appendChild(commentDiv);
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

export function fetchLikedPosts(body) {
    fetch(`/api/liked-posts`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            const posts = data.LikedPosts;
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

                     if(data.LoggedIn) {
                        const likeButtons = makeElement('div', '', 'post-buttons', '', '')
                        
                        likeButtons.innerHTML = `
                            <button class="like-button" data-post-id="${post.ID}">Like</button>
                            <button class="dislike-button" data-post-id="${post.ID}">Dislike</button>
                            `
                        likeButtons.querySelector('.like-button').addEventListener('click', () => handleLikeDislike(post.ID, true));
                        likeButtons.querySelector('.dislike-button').addEventListener('click', () => handleLikeDislike(post.ID, false));
                        
                        postBox.appendChild(likeButtons);
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
                                if(data.LoggedIn) {
                                    const commentButtons = makeElement('div', '', 'comment-buttons', '', '')
                                    commentButtons.innerHTML = `
                                    <button class="like-button" data-comment-id="${comment.ID}">Like</button>
                                    <button class="dislike-button" data-comment-id="${comment.ID}">Dislike</button>
                                    `
                                    commentButtons.querySelector('.like-button').addEventListener('click', () => handleCommentLikes(comment.ID, true));
                                    commentButtons.querySelector('.dislike-button').addEventListener('click', () => handleCommentLikes(comment.ID, false));

                                    commentDiv.appendChild(commentButtons)
                                }
                            comments.appendChild(commentDiv);
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
