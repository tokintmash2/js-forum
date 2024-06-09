import { renderHomePage } from "./render.js";

export function handleLikeDislike(postId, isLike) {
    const url = isLike ? '/api/likePost' : '/api/dislikePost';
    console.log('Post ID:', postId);
    console.log('URL:', url);
    
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ post_id: postId })
    })
    .then(response => {
        console.log('Response Status:', response.status);

        if (!response.ok) {
            return response.text().then(text => { throw new Error(text) });
        }
        return response.json();
    })
    
    .then(data => {
        if (data.success) {
            console.log(data.newLikeCount); // New like count
            const postDetailsElement = document.querySelector(`.like-button[data-post-id="${postId}"]`)
                .closest('.post').querySelector('.post-details');
            if (postDetailsElement) {
                const likesDislikesElement = Array.from(postDetailsElement.querySelectorAll('p'))
                    .find(p => p.textContent.includes('Likes:'));
                if (likesDislikesElement) {
                    // Extract dislikes count from the current text
                    const dislikesCount = likesDislikesElement.textContent.match(/Dislikes: (\d+)/)[1];
                    likesDislikesElement.textContent = `
                        Likes: ${data.newLikeCount}, Dislikes: ${data.newDislikeCount}`;
                } else {
                    console.error('Likes/Dislikes element not found.');
                }
            } else {
                console.error('Post details element not found.');
            }
            }
        })
    
    .catch(error => {
        console.error('Error liking/disliking post:', error);
    });
}

export function handleCommentLikes(commentId, isLike) {
    const url = isLike ? '/api/likeComment' : '/api/dislikeComment';
    console.log('Comment ID:', commentId);
    console.log('URL:', url);
    
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ comment_id: commentId })
    })
    .then(response => {
        console.log('Response Status:', response.status);

        if (!response.ok) {
            return response.text().then(text => { throw new Error(text) });
        }
        return response.json();
    })
    
    .then(data => {
        if (data.success) {
            console.log(data.newLikeCount); // New like count
            const postDetailsElement = document.querySelector(`.like-button[data-comment-id="${commentId}"]`)
                .closest('.comment');
            if (postDetailsElement) {
                const likesDislikesElement = Array.from(postDetailsElement.querySelectorAll('p'))
                    .find(p => p.textContent.includes('Likes:'));
                if (likesDislikesElement) {
                    // Extract dislikes count from the current text
                    const dislikesCount = likesDislikesElement.textContent.match(/Dislikes: (\d+)/)[1];
                    likesDislikesElement.textContent = `
                        Likes: ${data.newLikeCount}, Dislikes: ${data.newDislikeCount}`;
                } else {
                    console.error('Likes/Dislikes element not found.');
                }
            } else {
                console.error('Post details element not found.');
            }
            }
        })
    
    .catch(error => {
        console.error('Error liking/disliking post:', error);
    });
}

export function handleLogout() {
    fetch('/api/logout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            // Perform any client-side clean-up, such as updating the UI
            console.log('Logout successful');
            // Redirect or render the home page
            // renderHomePage();
        } else {
            console.error('Logout failed');
        }
    })
    .catch(error => {
        console.error('Error during logout:', error);
    });
}

export function creatPostHandler() {
    const title = document.getElementById('title').value;
    const content = document.getElementById('content').value;
    const categoryIDs = Array.from(document.querySelectorAll('input[name="categoryIDs"]:checked')).map(checkbox => parseInt(checkbox.value));

    const postData = {
        title: title,
        content: content,
        categoryIDs: categoryIDs
    };

    fetch('/api/create-post', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(postData)
    })
    .then(response => {
        if (!response.ok) {
            return response.text().then(text => { throw new Error(text) });
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            renderHomePage()
            console.log('Post created successfully');
            // Redirect or update the UI as needed
        } else {
            console.error('Error creating post');
        }
    })
    .catch(error => {
        console.error('Error creating post:', error);
    });
}

export function createCommentHandler(post) {
   
    const postID = parseInt(post)
    // const postID = document.getElementById('content').value
    const content = document.getElementsByName('comment_content').value;
    // const categoryIDs = Array.from(document.querySelectorAll('input[name="categoryIDs"]:checked')).map(checkbox => parseInt(checkbox.value));

    const commentData = {
        post_id: postID,
        content: content,
        // categoryIDs: categoryIDs
    };

    fetch('/api/add-comment', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(commentData)
    })
    .then(response => {
        if (!response.ok) {
            return response.text().then(text => { throw new Error(text) });
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            renderHomePage()
            console.log('Comment added successfully');
            // Redirect or update the UI as needed
        } else {
            console.error('Error adding comment');
        }
    })
    .catch(error => {
        console.error('Error adding comment:', error);
    });
}