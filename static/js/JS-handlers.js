import { renderHomePage } from "./render.js";
import { appendNewComment } from "./append.js";

export function login() {
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    if (!email || !password) {
        alert('All fields are required.')
        return
    }

    console.log("u:", email, "p:", password)

    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email: email, password: password })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            console.log('Login successful');
            connectWebSocket()
            document.getElementById('app').innerHTML = ''
            renderHomePage()

        } else {
            alert('Login failed: ' + data.message);
        }
    })
    .catch(error => console.error('Error:', error));
}

function connectWebSocket() {
    const socket = new WebSocket('ws://localhost:4000/ws');

    socket.onopen = function() {
        console.log('WebSocket connection established');
    };

    socket.onclose = function(event) {
        if (event.wasClean) {
            console.log('WebSocket connection closed cleanly');
        } else {
            console.log('WebSocket connection closed unexpectedly');
        }
    };

    socket.onerror = function(error) {
        console.error('WebSocket error:', error);
    };
}

export function signup() {
    const username = document.getElementById('signup-username').value;
    const email = document.getElementById('signup-email').value;
    const password = document.getElementById('signup-password').value;

    if (!username || !email || !password) {
        alert('All fields are required.')
        return
    }

    if(!validateEmail(email)) {
        alert('Please use a valid email address!')
        return
    }

    fetch('/api/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username: username, email: email, password: password })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            alert('Signup successful');
            // Handle successful signup (e.g., redirect to the login page)
        } else {
            alert('Signup failed: ' + data.message);
        }
    })
    .catch(error => console.error('Error:', error));
}

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
            document.getElementById('app').innerHTML = ''
            console.log('Logout successful');
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
        } else {
            console.error('Error creating post');
        }
    })
    .catch(error => {
        console.error('Error creating post:', error);
    });
}

export function createCommentHandler(post_ID) {
    const postID = parseInt(post_ID)
    const content = document.getElementById(`comment-post-id-${postID}`).value;

    console.log(postID, content)

    const commentData = {
        post_id: postID,
        content: content,
    };

    console.log('Comment Data:', commentData);

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
            console.log('data:', data);
            console.log('Comment added successfully');
            appendNewComment(data.comment);
            document.getElementById(`comment-post-id-${postID}`).value = ''; // Clear input
        } else {
            console.error('Error adding comment');
        }
    })
    .catch(error => {
        console.error('Error adding comment:', error);
    });
}

function validateEmail(email) {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(String(email).toLowerCase());
}

