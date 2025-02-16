document.querySelector(".login-form").addEventListener("submit", async function (e) {
    e.preventDefault();
    console.log("Form submitted"); // Проверка отправки формы

    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    try {
        const response = await fetch("/api/auth/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password })
        });

        if (response.ok) {
            console.log("Login successful");
            const data = await response.json(); // Получаем данные из ответа

            // Сохраняем имя пользователя и email в localStorage
            localStorage.setItem("username", username);
            localStorage.setItem("token", data.token); // Сохраняем токен

            //Проверяем, пришел ли email
            if (data.email) {
                localStorage.setItem("email", data.email);
            }

            // Перенаправление на главную страницу
            window.location.href = "/FrontEnd/public/index.html";
        } else {
            const errorData = await response.json();
            displayError(errorData.error || "Login failed");
        }
    } catch (error) {
        console.error("Error during login:", error);
        displayError("An error occurred. Please try again.");
    }
});

// Функция для отображения ошибки
function displayError(message) {
    let errorElement = document.querySelector(".error-message");
    if (!errorElement) {
        errorElement = document.createElement("div");
        errorElement.className = "error-message";
        document.querySelector(".login-container").appendChild(errorElement);
    }
    errorElement.textContent = message;
}