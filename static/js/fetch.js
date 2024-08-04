import { makeElements} from "./make-elements.js";
import { handleLikeDislike, createCommentHandler } from "./JS-handlers.js"
import { reinitializeCommentFunctionality } from "./commentToggle.js"
import { appendComments, appendOnlineUsers } from "./append.js";

export function fetchOnlineUsers() {
    fetch('/api/online-user')
        .then(response => response.json())
        .then(users => appendOnlineUsers(users)) // Append users to sidebar
        .catch(error => {
            console.error('Error fetching online users:', error);
        });
}

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
    const mainContent = body.querySelector('.main-content');
    // const chatSection = body.children[1]
   return fetch(`/api/category/${id}`) // Promise so, render function can wait
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
                    const postDiv = makeElements({type:'div', name: `${post.ID}`, classNames:'post'});
    
                    const postBox = makeElements({type:'div', classNames:'content-box', contents:`
                        <h3>${post.Title}</h3><p>${post.Content}</p>`});
                    postDiv.appendChild(postBox);
    
                    const postDetails = makeElements({type:'div', classNames:'post-details', contents:`
                        <p>Author: ${post.Author}</p>
                        <p>Likes: ${post.Likes}, Dislikes: ${post.Dislikes}</p>
                        <p>Created at: ${new Date(post.CreatedAt).toLocaleDateString()}</p>`});
                    postBox.appendChild(postDetails);

                     if(data.LoggedIn) {

                        reinitializeCommentFunctionality()

                        const likeButtons = makeElements({type:'div', classNames:'post-buttons'})
                        const addCommentDiv = makeElements({type:'div', name:'add-comment', classNames:'add-comment-section'})
                        
                        likeButtons.innerHTML = `
                            <button class="like-button" data-post-id="${post.ID}">Like</button>
                            <button class="dislike-button" data-post-id="${post.ID}">Dislike</button>
                            `
                        likeButtons.querySelector('.like-button').addEventListener('click', () => handleLikeDislike(post.ID, true));
                        likeButtons.querySelector('.dislike-button').addEventListener('click', () => handleLikeDislike(post.ID, false));

                        addCommentDiv.innerHTML = `
                            <button class="add-comment-btn">Write Comment</button>
                            <textarea name="comment_content" id="comment-post-id-${post.ID}" rows="6" cols="50" placeholder="Leave your comment here"></textarea>
                            <button class="sumbit-comment-btn" type="submit">Add Comment</button>
                            `
                        addCommentDiv.querySelector('.sumbit-comment-btn').addEventListener('click', () => {
                            console.log(post.ID)
                            createCommentHandler(post.ID)
                        })                        

                        postBox.appendChild(likeButtons);
                        postBox.appendChild(addCommentDiv)
                     }
                    
                    // Comments
                    const commentBox = makeElements({type:'div', classNames:'comment-container'});
                    const comments = makeElements({type:'div', classNames: 'comments'});
                    
                    appendComments(post, data, comments)

                    commentBox.appendChild(comments);
                    postBox.appendChild(commentBox);
                    mainContent.appendChild(postDiv)                    
                });
            } else {
                const noPosts = makeElements('div', '', '', 'No posts in this category.', '')
                body.appendChild(noPosts);
            }
            fetchOnlineUsers()
            body.appendChild(mainContent);
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
                const postDiv = makeElements({type:'div', name:'JSPost', classNames:'post'})

                const postBox = makeElements({type:'div', classNames:'content-box', contents:`
                    <h3>${post.Title}</h3><p>${post.Content}</p>`})
                postDiv.appendChild(postBox);

                const postDetails = makeElements({type:'div', classNames:'post-details', contents:`
                    <p>Author: ${post.Author}</p>
                    <p>Likes: ${post.Likes}, Dislikes: ${post.Dislikes}</p>
                    <p>Created at: ${new Date(post.CreatedAt).toLocaleDateString()}</p>`})
                postBox.appendChild(postDetails)

                // Comments
                const commentBox = makeElements({type:'div', classNames:'comment-container'})
                const comments = makeElements({type:'div', classNames:'comments'})

                if (Array.isArray(post.Comments)) {
                    post.Comments.forEach(comment => {
                        const commentDiv = makeElements({type: 'div', classNames:'comment', contents:`
                        <p>${comment.Content}</p>
                        <p>Author: ${comment.Author}</p>
                        <p>Likes: ${comment.Likes}, Dislikes: ${comment.Dislikes}</p>`})
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
