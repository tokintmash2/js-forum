import { makeElements } from './make-elements.js';
import { handleCommentLikes } from './JS-handlers.js';
import { loadConversation } from './fetch.js';

export function appendOnlineUsers(users) {
    const sidebar = document.getElementsByClassName('sidebar')[0]
    sidebar.innerHTML = ''
    const onlineUsers = makeElements({type: 'div', name: 'online-users'}) // Div for users
        if(users) {
            users.forEach(user => {
                const onlineUser = makeElements({
                    type: 'div',
                    name: 'user',
                    classNames: `${user.ID}`,
                    contents: `${user.Username}`
                }) // User
                onlineUser.addEventListener('click', () => loadConversation(user.ID, user.Username));
                onlineUsers.appendChild(onlineUser)
            })    
        } else {
            const onlineUser = makeElements({type: 'div', name: 'user', contents: `No online users`})
            onlineUsers.appendChild(onlineUser)
        }    
    sidebar.appendChild(onlineUsers)
    appendChat()
}

export function appendChat() {
    const sidebar = document.getElementsByClassName('sidebar')[0]
    const chatArea = makeElements({type: 'div', classNames: 'chatArea'})
    const chat = makeElements({type: 'div', classNames: 'chat'})

    chatArea.appendChild(chat);    

    const chatInput = makeElements({ type: 'textarea', classNames: 'chat-input', placeholder: 'Write a message...' });
    const sendButton = makeElements({ type: 'button', classNames: 'sendButton', contents: 'Send' });
    
    chatArea.appendChild(chatInput);
    chatArea.appendChild(sendButton);
    sidebar.appendChild(chatArea)
}



export function appendNewComment(comment) {

    const post_ID = comment.postID

    const postDiv = document.getElementById(post_ID)
    const commentBox = postDiv.querySelector('.comment-container');
    const comments = postDiv.querySelector('.comments');

    const commentDiv = document.createElement('div');

    const likeButtons = makeElements({ type: 'div', classNames: 'comment-buttons' })

    commentDiv.className = 'comment';
    commentDiv.innerHTML = `
        <p>${comment.content}</p>
        <p>Author: ${comment.author}</p>
        <p>Likes: ${comment.likes}, Dislikes: ${comment.dislikes}</p>
    `;

    likeButtons.innerHTML = `
            <button class="like-button" data-comment-id="${comment.ID}">Like</button>
            <button class="dislike-button" data-comment-id="${comment.ID}">Dislike</button>
            `
    likeButtons.querySelector('.like-button').addEventListener('click', () => handleCommentLikes(comment.ID, true));
    likeButtons.querySelector('.dislike-button').addEventListener('click', () => handleCommentLikes(comment.ID, false));

    commentDiv.appendChild(likeButtons)
    comments.appendChild(commentDiv);
    commentBox.appendChild(comments);
    postDiv.appendChild(commentBox);
}

export function appendComments(post, data, commentsContainer) {
    if (Array.isArray(post.Comments)) {
        post.Comments.forEach(comment => {
            const commentDiv = makeElements({
                type: 'div',
                classNames: 'comment',
                contents: `
                    <p>${comment.Content}</p>
                    <p>Author: ${comment.Author}</p>
                    <p>Likes: ${comment.Likes}, Dislikes: ${comment.Dislikes}</p>`
            });
            if (data.LoggedIn) {
                const commentButtons = makeElements({ type: 'div', classNames: 'comment-buttons' });
                commentButtons.innerHTML = `
                    <button class="like-button" data-comment-id="${comment.ID}">Like</button>
                    <button class="dislike-button" data-comment-id="${comment.ID}">Dislike</button>
                `;
                commentButtons.querySelector('.like-button').addEventListener('click', () => handleCommentLikes(comment.ID, true));
                commentButtons.querySelector('.dislike-button').addEventListener('click', () => handleCommentLikes(comment.ID, false));

                commentDiv.appendChild(commentButtons);
            }
            commentsContainer.appendChild(commentDiv);
        });
    }
}