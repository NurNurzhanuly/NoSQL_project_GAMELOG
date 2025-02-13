document.addEventListener('DOMContentLoaded', function() {
    const username = localStorage.getItem('username');
    const authMessage = document.getElementById('auth-message');
    const cartContent = document.querySelector('.cart-content');

    if (username) {
        // Пользователь авторизован
        cartContent.style.display = 'block';
    } else {
        // Пользователь не авторизован
        authMessage.style.display = 'block';
    }
});