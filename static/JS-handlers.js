export function handleLikeDislike(postId, isLike) {
    const url = isLike ? '/api/likePost' : '/api/dislikePost';
    console.log(postId)
    console.log(url)
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ post_id: postId })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            // Update the like count dynamically
            const likeCountElement = document.querySelector(`.like-button[data-post-id="${postId}"]`).closest('.post').querySelector('.like-count');
            likeCountElement.textContent = data.newLikeCount;
        }
    })
    .catch(error => {
        console.error('Error liking/disliking post:', error);
    });
}