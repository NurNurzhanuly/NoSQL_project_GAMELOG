document.addEventListener('DOMContentLoaded', function () {
    const username = localStorage.getItem('username');
    const authMessage = document.getElementById('auth-message');
    const cartContent = document.querySelector('.cart-content');

    console.log('Username in localStorage:', localStorage.getItem('username'));
    console.log('Token in localStorage:', localStorage.getItem('token'));

    if (username) {
        // Пользователь авторизован
        cartContent.style.display = 'block';
        loadCartItems();
    } else {
        // Пользователь не авторизован
        authMessage.style.display = 'block';
    }

    async function loadCartItems() {
        const cartItemsContainer = document.querySelector('.cart-items');

        try {
            console.log('Fetching cart items...'); // Лог перед запросом

            const response = await fetch('/api/cart', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (!response.ok) {
                console.error('Failed to fetch cart items. Status:', response.status);
                throw new Error('Failed to fetch cart items');
            }

            const cartData = await response.json();
            console.log('Cart data received:', cartData); // Лог данных корзины

            if (!cartData.items || cartData.items.length === 0) {
                cartItemsContainer.innerHTML = '<p>Your cart is empty.</p>';
                return;
            }

            cartItemsContainer.innerHTML = ''; // Очистить контейнер перед добавлением элементов

            let cartTotal = 0; // Общая сумма

            // Перебираем все элементы корзины и отображаем их
            for (const item of cartData.items) {
                console.log('Fetching game details for item:', item);

                const gameResponse = await fetch(`/api/games/${item.game_id}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                });

                if (!gameResponse.ok) {
                    console.error('Failed to fetch game details. Status:', gameResponse.status);
                    throw new Error('Failed to fetch game details');
                }

                const gameData = await gameResponse.json();
                console.log('Game data received:', gameData);

                // Отображение данных игры
                const cartItemDiv = document.createElement('div');
                cartItemDiv.classList.add('cart-item');

                cartItemDiv.innerHTML = `
                <img src="${gameData.image}" alt="${gameData.title}">
                <div class="cart-item-details">
                    <h3 class="cart-item-title">${gameData.title}</h3>
                    <p class="cart-item-price">$${gameData.price.toFixed(2)}</p>
                </div>
                <button class="remove-from-cart-btn" data-id="${item.game_id}">Remove</button>
            `;

                cartItemsContainer.appendChild(cartItemDiv);
                cartTotal += gameData.price;
            }

            // Обновляем общую стоимость корзины
            document.querySelector('.cart-total p').textContent = `Total: $${cartTotal.toFixed(2)}`;

            // Добавляем обработчики удаления товаров из корзины
            const removeFromCartButtons = document.querySelectorAll('.remove-from-cart-btn');
            removeFromCartButtons.forEach(button => {
                button.addEventListener('click', async function () {
                    const gameID = button.dataset.id;

                    try {
                        const deleteResponse = await fetch(`/api/cart/${gameID}`, {
                            method: 'DELETE',
                            headers: {
                                'Content-Type': 'application/json',
                                'Authorization': `Bearer ${localStorage.getItem('token')}`
                            }
                        });

                        if (!deleteResponse.ok) {
                            throw new Error('Failed to remove item from cart');
                        }

                        // Удаляем элемент из DOM после успешного удаления
                        button.parentElement.remove();

                        // Пересчитываем общую стоимость
                        cartTotal -= parseFloat(button.parentElement.querySelector('.cart-item-price').textContent.replace('$', ''));
                        document.querySelector('.cart-total p').textContent = `Total: $${cartTotal.toFixed(2)}`;
                    } catch (error) {
                        console.error('Error removing item from cart:', error);
                    }
                });
            });
        } catch (error) {
            console.error('Error loading cart items:', error);
            cartItemsContainer.innerHTML = '<p>Failed to load cart items. Please try again later.</p>';
        }
    }
});
