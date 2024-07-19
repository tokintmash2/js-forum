// JavaScript code to handle comment toggling
// An event listener to the document that captures clicks on "Write Comment" buttons
document.addEventListener('DOMContentLoaded', function() {
    const isNewPostCreated = localStorage.getItem('isNewPostCreated');
    
    if (isNewPostCreated === 'true') {
        // Reinitialize the JavaScript functionality for handling "Write Comment" buttons
        reinitializeCommentFunctionality();
        
        // After reinitialization, remove the flag from localStorage
        localStorage.removeItem('isNewPostCreated');
    }
});

// Function to reinitialize the JavaScript functionality for comments
export function reinitializeCommentFunctionality() {
    // Find all the "Write Comment" buttons
    const commentButtons = document.querySelectorAll('.add-comment-btn');

    // Loop through each button and reattach event listeners
    commentButtons.forEach(button => {
        button.addEventListener('click', function(event) {
            // Toggle the comment form display
            const form = event.target.nextElementSibling;
            form.style.display = (form.style.display === 'none') ? 'block' : 'none';
        });
    });
}

// After creating a new post and redirecting or refreshing the page
// Set the flag in localStorage to indicate a new post was created
localStorage.setItem('isNewPostCreated', 'true');

function toggleCommentForm(button) {
    const commentForm = button.nextElementSibling;
    if (commentForm.style.display === 'none' || commentForm.style.display === '') {
        commentForm.style.display = 'block';
    } else {
        commentForm.style.display = 'none';
    }
}