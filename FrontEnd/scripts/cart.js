document.addEventListener('DOMContentLoaded', function () {
    const username = localStorage.getItem('username');
    const authMessage = document.getElementById('auth-message');
    const cartContent = document.querySelector('.cart-content');
    const cartItemsContainer = document.querySelector('.cart-items');
    const buyAllBtn = document.getElementById("buy-all-btn");
    const confirmPurchaseBtn = document.getElementById("confirm-purchase-btn");
    const cancelPurchaseBtn = document.getElementById("cancel-purchase-btn");
    const buyConfirmationModal = document.getElementById("buy-confirmation-modal");
    const totalPriceMessage = document.getElementById("total-price-message");
    let cartTotal = 0;

    if (username) {
        cartContent.style.display = 'block';
        loadCartItems();
    } else {
        authMessage.style.display = 'block';
    }

    async function loadCartItems() {
        try {
            const response = await fetch('/api/cart', getAuthHeaders());
            if (!response.ok) throw new Error('Failed to fetch cart items');

            const cartData = await response.json();
            updateCartDisplay(cartData.items);
        } catch (error) {
            cartItemsContainer.innerHTML = '<p>Failed to load cart items. Please try again later.</p>';
            console.error(error);
        }
    }

    function getAuthHeaders() {
        return {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        };
    }

    async function fetchGameDetails(gameId) {
        const response = await fetch(`/api/games/${gameId}`, getAuthHeaders());
        if (!response.ok) throw new Error('Failed to fetch game details');
        return await response.json();
    }

    async function updateCartDisplay(items) {
        cartItemsContainer.innerHTML = items.length ? '' : '<p>Your cart is empty.</p>';
        cartTotal = 0;

        for (const item of items) {
            try {
                const gameData = await fetchGameDetails(item.game_id);
                cartTotal += gameData.price;
                cartItemsContainer.appendChild(createCartItemElement(item.game_id, gameData));
            } catch (error) {
                console.error(error);
            }
        }
        updateCartTotal();
    }

    function createCartItemElement(gameId, gameData) {
        const cartItemDiv = document.createElement('div');
        cartItemDiv.classList.add('cart-item');
        cartItemDiv.innerHTML = `
            <img src="${gameData.image}" alt="${gameData.title}">
            <div class="cart-item-details">
                <h3 class="cart-item-title">${gameData.title}</h3>
                <p class="cart-item-price">$${gameData.price.toFixed(2)}</p>
            </div>
            <button class="remove-from-cart-btn" data-id="${gameId}">Remove</button>
        `;
        cartItemDiv.querySelector('.remove-from-cart-btn').addEventListener('click', () => removeFromCart(gameId, cartItemDiv, gameData.price));
        return cartItemDiv;
    }

    async function removeFromCart(gameId, itemElement, price) {
        try {
            const response = await fetch(`/api/cart/${gameId}`, {
                method: 'DELETE',
                headers: getAuthHeaders().headers
            });
            if (!response.ok) throw new Error('Failed to remove item from cart');

            itemElement.remove();
            cartTotal -= price;
            updateCartTotal();
        } catch (error) {
            console.error('Error removing item from cart:', error);
        }
    }

    function updateCartTotal() {
        document.querySelector('.cart-total p').textContent = `Total: $${cartTotal.toFixed(2)}`;
    }

    buyAllBtn.addEventListener("click", function () {
        console.log("Buy All button clicked!"); // Проверка

        totalPriceMessage.textContent = `Do you want to buy these games for a total of $${cartTotal.toFixed(2)}?`;
        buyConfirmationModal.classList.remove("hidden");
    });

    confirmPurchaseBtn.addEventListener("click", async function () {
        try {
            const response = await fetch("/api/cart/purchase", {
                method: "POST",
                headers: { "Authorization": `Bearer ${localStorage.getItem("token")}` }
            });
            const data = await response.json();
            alert(data.message || "Purchase completed successfully!");
            buyConfirmationModal.classList.add("hidden");
            loadCartItems();
        } catch (error) {
            alert("Error completing purchase");
            console.error(error);
        }
    });

    cancelPurchaseBtn.addEventListener("click", function () {
        buyConfirmationModal.classList.add("hidden");
    });
});
