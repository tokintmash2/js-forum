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

// Post details element structure:
    // <div class="post-details">
    // <p>Author: Papakoi</p>
    // <p>Likes: 2, Dislikes: 1</p>
    // <p>Created at: 12/22/2023</p></div>
    // Must update Likes/Dislikes 