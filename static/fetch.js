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
                const postDiv = document.createElement('div')
                postDiv.id = 'JSPost'
                postDiv.classList.add('post')
                const postBox = document.createElement('div')
                postBox.classList.add('content-box');
                postBox.innerHTML = `<h3>${post.Title}</h3><p>${post.Content}</p>` 
                postDiv.appendChild(postBox);
                const postDetails = document.createElement('div')
                postDetails.classList.add('post-details')
                postDetails.innerHTML = `
                    <p>Author: ${post.Author}</p>
                    <p>Likes: ${post.Likes}, Dislikes: ${post.Dislikes}</p>
                    <p>Created at: ${new Date(post.CreatedAt).toLocaleDateString()}</p>` 
                postBox.appendChild(postDetails)
                const commentBox = document.createElement('div')
                commentBox.classList.add('comment-container')
                const comments = document.createElement('div')
                comments.classList.add('comments')
                if (Array.isArray(post.Comments)) {
                    post.Comments.forEach(comment => {
                    const commentDiv = document.createElement('div')
                    commentDiv.classList.add('comment')
                    commentDiv.innerHTML = `
                    <p>${comment.Content}</p>
                    <p>Author: ${comment.Author}</p>
                    <p>Likes: ${comment.Likes}, Dislikes: ${comment.Dislikes}</p>` 
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
            resolve (userData)
        })
        .catch(error => {
            console.error('Error loading the loginData:', error);
            throw error
        });
    })    
}