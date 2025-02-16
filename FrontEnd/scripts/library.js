document.addEventListener('DOMContentLoaded', function () {
    const username = localStorage.getItem('username');
    const authMessage = document.getElementById('auth-message');
    const accountInfo = document.getElementById('account-info');
    const libraryHeader = document.getElementById('library-header');
    const libraryList = document.getElementById('library-list');

    // check if user is logged in
    if (username) {
        accountInfo.style.display = 'block';
        libraryHeader.style.display = 'block';
        libraryList.style.display = 'grid';
    } else {
        authMessage.style.display = 'block';
    }
});