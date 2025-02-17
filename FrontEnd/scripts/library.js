document.addEventListener('DOMContentLoaded', async function () {
    const username = localStorage.getItem('username');
    const authMessage = document.getElementById('auth-message');
    const accountInfo = document.getElementById('account-info');
    const libraryHeader = document.getElementById('library-header');
    const libraryList = document.getElementById('library-list');

    if (username) {
        accountInfo.style.display = 'block';
        libraryHeader.style.display = 'block';
        libraryList.style.display = 'grid';

        try {
            const response = await fetch('/api/user/library', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (response.ok) {
                const games = await response.json();
                displayGames(games);
            } else {
                console.error('Failed to fetch library games');
            }
        } catch (error) {
            console.error('Error fetching library:', error);
        }
    } else {
        authMessage.style.display = 'block';
    }
});

function displayGames(games) {
    const libraryList = document.getElementById('library-list');
    libraryList.innerHTML = '';

    games.forEach(game => {
        const gameCard = document.createElement('div');
        gameCard.classList.add('game-card');

        gameCard.innerHTML = `
            <img src="${game.imageUrl}" alt="${game.title}" class="game-image">
            <h3>${game.title}</h3>
            <p>${game.description}</p>
        `;

        libraryList.appendChild(gameCard);
    });
}
