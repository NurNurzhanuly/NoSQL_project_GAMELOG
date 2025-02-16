document.addEventListener("DOMContentLoaded", async function () {
    const gamesGrid = document.getElementById("all-games-grid");
    const searchInput = document.querySelector('.search-input');

    // Функция загрузки игр
    async function loadGames(searchTerm = "") {
        try {
            const url = searchTerm
                ? `/api/games/search?title=${encodeURIComponent(searchTerm)}`
                : "/api/games/getall";
            const response = await fetch(url);
            if (!response.ok) throw new Error("Failed to fetch games");
            const games = await response.json();
            gamesGrid.innerHTML = ""; // Очищаем текущий контент
            games.forEach((game) => {
                const gameCard = document.createElement("div");
                gameCard.classList.add("game-card");
                gameCard.innerHTML = `
                    <img src="${game.image_url}" alt="${game.title}">
                    <div class="game-info">
                        <h3>${game.title}</h3>
                        <p>${game.description}</p>
                        <p>Genre: ${game.genre}</p>
                        <p>Developer: ${game.developer}</p>
                        <p>Price: $${game.price.toFixed(2)}</p>
                        <button class="add-to-cart-btn" data-game-id="${game.id}">Add to Cart</button>
                    </div>
                `;
                gamesGrid.appendChild(gameCard);
            });

            // Добавляем обработчики событий на новые кнопки
            addCartButtonListeners();
        } catch (error) {
            console.error("Failed to load games:", error);
        }
    }

    // Функция добавления обработчиков событий для кнопок
    function addCartButtonListeners() {
        const isUserLoggedIn = !!localStorage.getItem("token");
        document.querySelectorAll(".add-to-cart-btn").forEach((button) => {
            button.addEventListener("click", async function () {
                if (!isUserLoggedIn) {
                    window.location.href = "/FrontEnd/public/login.html";
                    return;
                }

                const gameId = this.getAttribute("data-game-id");
                try {
                    const token = localStorage.getItem("token"); // Получаем токен из localStorage
                    const response = await fetch(`/api/cart`, {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                            "Authorization": `Bearer ${token}`, // Добавляем токен в заголовок
                        },
                        body: JSON.stringify({
                            game_id: gameId,
                        }),
                    });
                    if (response.ok) {
                        alert("Game added to cart!");
                    } else {
                        throw new Error("Failed to add to cart");
                    }
                } catch (error) {
                    console.error("Error:", error);
                    alert("An error occurred while adding to cart");
                }
            });
        });
    }

    // Загружаем все игры по умолчанию
    loadGames();

    // Добавляем обработчик событий для поиска
    searchInput.addEventListener("input", function () {
        const searchTerm = searchInput.value.trim();
        loadGames(searchTerm);
    });
});
