document.addEventListener('DOMContentLoaded', function () {
    const username = localStorage.getItem('username');
    const authMessage = document.getElementById('auth-message');
    const cartContent = document.querySelector('.cart-content');

    if (username) {
        cartContent.style.display = 'block';
        loadCartItems();
    } else {
        authMessage.style.display = 'block';
    }

    async function loadCartItems() {
        const cartItemsContainer = document.querySelector('.cart-items');

        try {
            const response = await fetch('/api/cart', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to fetch cart items');
            }

            const cartData = await response.json();

            if (!cartData.items || cartData.items.length === 0) {
                cartItemsContainer.innerHTML = '<p>Your cart is empty.</p>';
                return;
            }

            cartItemsContainer.innerHTML = '';
            let cartTotal = 0;

            for (const item of cartData.items) {
                const gameResponse = await fetch(`/api/games/${item.game_id}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                });

                if (!gameResponse.ok) {
                    throw new Error('Failed to fetch game details');
                }

                const gameData = await gameResponse.json();
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

            document.querySelector('.cart-total p').textContent = `Total: $${cartTotal.toFixed(2)}`;

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

                        button.parentElement.remove();
                        cartTotal -= parseFloat(button.parentElement.querySelector('.cart-item-price').textContent.replace('$', ''));
                        document.querySelector('.cart-total p').textContent = `Total: $${cartTotal.toFixed(2)}`;
                    } catch (error) {
                        console.error('Error removing item from cart:', error);
                    }
                });
            });
        } catch (error) {
            cartItemsContainer.innerHTML = '<p>Failed to load cart items. Please try again later.</p>';
        }
    }

    const buyAllBtn = document.getElementById("buy-all-btn");
    const confirmPurchaseBtn = document.getElementById("confirm-purchase-btn");
    const cancelPurchaseBtn = document.getElementById("cancel-purchase-btn");
    const buyConfirmationModal = document.getElementById("buy-confirmation-modal");
    const totalPriceMessage = document.getElementById("total-price-message");

    let totalPrice = 0.00;

    // Открытие окна подтверждения покупки
    buyAllBtn.addEventListener("click", function () {
        totalPriceMessage.textContent = `Do you want to buy these games for a total of $${totalPrice.toFixed(2)}?`;
        buyConfirmationModal.classList.remove("hidden");
    });

    // Подтверждение покупки
    confirmPurchaseBtn.addEventListener("click", function () {
        fetch("/api/cart/purchase", {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`,
            },
        })
            .then(response => response.json())
            .then(data => {
                alert(data.message || "Purchase completed successfully!");
                buyConfirmationModal.classList.add("hidden");
                loadCartItems(); // обновляем корзину
            })
            .catch(error => {
                alert("Error completing purchase");
                console.error(error);
            });
    });

    // Отмена покупки
    cancelPurchaseBtn.addEventListener("click", function () {
        buyConfirmationModal.classList.add("hidden");
    });
});
