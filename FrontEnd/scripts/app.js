document.addEventListener('DOMContentLoaded', function() {
    console.log('Game Log frontend loaded successfully!');

    const menuToggle = document.getElementById('menuToggle');
    const sidebar = document.getElementById('sidebar');

    if (menuToggle && sidebar) {
         menuToggle.addEventListener('click', () => {
            sidebar.classList.toggle('open');
        });

        document.addEventListener('click', (event) => {
             if (!sidebar.contains(event.target) && !menuToggle.contains(event.target) && sidebar.classList.contains('open')) {
                sidebar.classList.remove('open');
             }
         });
    }

    const gameCardsHover = document.querySelectorAll('.game-card');
    gameCardsHover.forEach(card => {
        card.addEventListener('mouseenter', () => {
            card.style.transform = 'translateY(-10px)';
        });
        card.addEventListener('mouseleave', () => {
            card.style.transform = 'translateY(0)';
        });
    });

    // Search functionality
    const searchInput = document.querySelector('.search-input');
    if (searchInput) { // Проверяем, что searchInput существует
        const gameCardsSearch = document.querySelectorAll('.game-card');

        searchInput.addEventListener('input', function() {
            const searchTerm = searchInput.value.toLowerCase();

            gameCardsSearch.forEach(card => {
                const gameTitle = card.querySelector('h3').textContent.toLowerCase();
                if (gameTitle.includes(searchTerm)) {
                    card.style.display = 'flex';
                } else {
                    card.style.display = 'none';
                }
            });
        });
    }

    // Add to cart functionality
    const addToCartButtons = document.querySelectorAll('.add-to-cart-btn');
    addToCartButtons.forEach(button => {
        button.addEventListener('click', function() {
            const gameTitle = button.dataset.gameTitle;
            const gameImage = button.dataset.gameImage;

            // Get existing cart items from localStorage
            let cartItems = JSON.parse(localStorage.getItem('cartItems')) || [];

            // Add new item to cart
            cartItems.push({
                title: gameTitle,
                image: gameImage,
                quantity: 1 // You can adjust quantity later
            });

            // Save updated cart items to localStorage
            localStorage.setItem('cartItems', JSON.stringify(cartItems));

            alert(`${gameTitle} added to cart!`); // Optional feedback message
        });
    });

    // Cart page functionality
    const cartItemsContainer = document.querySelector('.cart-items');
    if (cartItemsContainer) { //* добавил проверку, чтобы скрипт выполнялся только на странице корзины*/
        let cartItems = JSON.parse(localStorage.getItem('cartItems')) || [];

        if (cartItems.length === 0) {
            cartItemsContainer.innerHTML = '<p>Your cart is empty.</p>';
            return;
        }

        function updateCartTotal() {
            let cartTotal = 0;
            cartItems.forEach(item => {
                cartTotal += 59.99 * item.quantity;
            });
            document.querySelector('.cart-total p').textContent = `Total: $${cartTotal.toFixed(2)}`;
        }

        cartItems.forEach((item, index) => {
            const cartItemDiv = document.createElement('div');
            cartItemDiv.classList.add('cart-item');

            cartItemDiv.innerHTML = `
                <img src="${item.image}" alt="${item.title}">
                <div class="cart-item-details">
                    <h3 class="cart-item-title">${item.title}</h3>
                    <p class="cart-item-price">$59.99</p>
                </div>
                <div class="cart-item-quantity">
                    <label for="quantity${index}">Quantity:</label>
                    <input type="number" id="quantity${index}" name="quantity${index}" value="${item.quantity}" min="1" data-index="${index}" class="quantity-input">
                </div>
                <button class="remove-from-cart-btn" data-index="${index}">Remove</button>
            `;

            cartItemsContainer.appendChild(cartItemDiv);
        });

        updateCartTotal();

        // Remove from cart functionality
        const removeFromCartButtons = document.querySelectorAll('.remove-from-cart-btn');
        removeFromCartButtons.forEach(button => {
            button.addEventListener('click', function() {
                const indexToRemove = parseInt(button.dataset.index);

                // Remove item from cartItems array
                cartItems.splice(indexToRemove, 1);

                // Save updated cart items to localStorage
                localStorage.setItem('cartItems', JSON.stringify(cartItems));

                // Reload cart page
                location.reload();
            });
        });

         // Quantity change functionality
        const quantityInputs = document.querySelectorAll('.cart-item-quantity input');
        quantityInputs.forEach(input => {
            input.addEventListener('change', function() {
                const indexToUpdate = parseInt(this.dataset.index);
                const newQuantity = parseInt(this.value);

                if (newQuantity > 0) {
                    cartItems[indexToUpdate].quantity = newQuantity;
                    localStorage.setItem('cartItems', JSON.stringify(cartItems));
                    updateCartTotal(); // Update the total on quantity change
                } else {
                    // If quantity is set to 0, remove the item
                    cartItems.splice(indexToUpdate, 1);
                    localStorage.setItem('cartItems', JSON.stringify(cartItems));
                    location.reload(); // Refresh the page to update the cart
                }
            });
        });
        // "Buy All" button functionality
        const buyAllButton = document.querySelector('.buy-all-button');
        buyAllButton.addEventListener('click', function() {
             if (cartItems.length > 0) {
                // Get existing library items from localStorage
                let libraryItems = JSON.parse(localStorage.getItem('libraryItems')) || [];

               // Add each item in the cart to the library
               cartItems.forEach(item => {
                   // Check if the game is already in the library
                   const gameExists = libraryItems.some(libraryItem => libraryItem.title === item.title);

                   if (!gameExists) {
                       libraryItems.push({
                           title: item.title,
                           image: item.image
                       });
                   }
               });

               // Save updated library items to localStorage
               localStorage.setItem('libraryItems', JSON.stringify(libraryItems));

               alert('Thank you for your purchase! Games added to your library.'); // Replace with actual checkout logic
               localStorage.removeItem('cartItems'); // Clear the cart
               location.reload(); // Refresh the page
           } else {
               alert('Your cart is empty.');
           }
        });
    }
     // Library page functionality
    const libraryList = document.getElementById('library-list');
    if (libraryList) {
        let libraryItems = JSON.parse(localStorage.getItem('libraryItems')) || [];
        const editLibraryBtn = document.getElementById('edit-library-btn');
        const deleteSelectedBtn = document.getElementById('delete-selected-btn');
        let isEditing = false; // Track editing mode

        if (libraryItems.length === 0) {
            libraryList.innerHTML = '<p>Your library is empty.</p>';
            return;
        }
         function updateLibraryDisplay() {
            libraryList.innerHTML = ''; // Clear existing list
            libraryItems.forEach(item => {
                const gameCard = document.createElement('div');
                gameCard.classList.add('game-card');
                const checkboxHTML = isEditing ? `<input type="checkbox" data-index="${index}" class="delete-checkbox">` : '';

                gameCard.innerHTML = `
                    <img src="${item.image}" alt="${item.title}">
                    <div class="game-card-content">
                        <h3 class="game-card-title">${item.title}</h3>
                          ${checkboxHTML}
                    </div>
                `;

            libraryList.appendChild(gameCard);
        });
    }

    updateLibraryDisplay();
        editLibraryBtn.addEventListener('click', function() {
        isEditing = !isEditing; // Toggle editing mode
        updateLibraryDisplay()

        if (isEditing) {
            editLibraryBtn.textContent = 'Exit Edit Mode';
            deleteSelectedBtn.style.display = 'block';
        } else {
            editLibraryBtn.textContent = 'Edit Library';
            deleteSelectedBtn.style.display = 'none';
        }
    });
            deleteSelectedBtn.addEventListener('click', function() {
        const checkedCheckboxes = document.querySelectorAll('.delete-checkbox:checked');
        const indicesToDelete = Array.from(checkedCheckboxes).map(checkbox => parseInt(checkbox.dataset.index)).sort((a, b) => b - a);

        indicesToDelete.forEach(index => {
            libraryItems.splice(index, 1);
        });

        localStorage.setItem('libraryItems', JSON.stringify(libraryItems));
        updateLibraryDisplay(); // Refresh the library display
        isEditing = false;
        editLibraryBtn.textContent = 'Edit Library';
        deleteSelectedBtn.style.display = 'none';
    });
}
});

document.addEventListener("DOMContentLoaded", function () {
    const username = localStorage.getItem("username");
    const navList = document.querySelector("#user-nav");

    if (username) {
        // Удалить кнопки "Login" и "Register"
        navList.innerHTML = "";

        // Добавить имя пользователя и кнопку выхода
        const userElement = document.createElement("li");
        const logoutElement = document.createElement("li");

        // Ссылка на библиотеку с именем пользователя
        const userLink = document.createElement("a");
        userLink.href = "/FrontEnd/public/library.html";
        userLink.textContent = username;

        // Кнопка выхода
        const logoutLink = document.createElement("a");
        logoutLink.href = "#";
        logoutLink.textContent = "Logout";
        logoutLink.addEventListener("click", function () {
            // Удаляем пользователя из localStorage и перезагружаем страницу
            localStorage.removeItem("username");
            window.location.href = "/FrontEnd/public/index.html";
        });

        userElement.appendChild(userLink);
        logoutElement.appendChild(logoutLink);
        navList.appendChild(userElement);
        navList.appendChild(logoutElement);
    }
});

// Account Information section
const usernameDisplay = document.getElementById('username-display');
const emailDisplay = document.getElementById('email-display');
const usernameInput = document.getElementById('username-input');
const emailInput = document.getElementById('email-input');
const editProfileBtn = document.getElementById('edit-profile-btn');
const saveProfileBtn = document.getElementById('save-profile-btn');

if (usernameDisplay && emailDisplay && usernameInput && emailInput && editProfileBtn && saveProfileBtn) {
    const username = localStorage.getItem('username');
    const email = localStorage.getItem('email');

    if (username) {
        usernameDisplay.textContent = username;
        usernameInput.value = username;
    }

    if (email) {
        emailDisplay.textContent = email;
        emailInput.value = email;
    }

    editProfileBtn.addEventListener('click', function() {
        usernameDisplay.style.display = 'none';
        emailDisplay.style.display = 'none';
        usernameInput.style.display = 'inline';
        emailInput.style.display = 'inline';
        editProfileBtn.style.display = 'none';
        saveProfileBtn.style.display = 'inline';
    });

    saveProfileBtn.addEventListener('click', async function() {
        const newUsername = usernameInput.value;
        const newEmail = emailInput.value;

        try {
            const response = await fetch('/api/user/update', { // Замените на ваш actual API endpoint
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + localStorage.getItem('token') // Если требуется авторизация
                },
                body: JSON.stringify({ username: newUsername, email: newEmail })
            });

            if (response.ok) {
                // Обработка успешного ответа
                console.log('Profile updated successfully!');
                localStorage.setItem('username', newUsername);

                // Проверяем, что email пришел с сервера
                if (newEmail) {
                    localStorage.setItem('email', newEmail);
                    emailDisplay.textContent = newEmail;
                }

                usernameDisplay.textContent = newUsername;

            } else {
                // Обработка ошибки
                console.error('Error updating profile:', response.statusText);
                alert('Failed to update profile. Please try again.'); // Или отобразите более информативное сообщение об ошибке
            }
        } catch (error) {
            console.error('Error updating profile:', error);
            alert('An error occurred. Please try again later.');
        }

        // Возвращаем все в исходное состояние
        usernameDisplay.style.display = 'inline';
        emailDisplay.style.display = 'inline';
        usernameInput.style.display = 'none';
        emailInput.style.display = 'none';
        editProfileBtn.style.display = 'inline';
        saveProfileBtn.style.display = 'none';
    });
}