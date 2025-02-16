document.querySelector(".login-form").addEventListener("submit", async function (e) {
    e.preventDefault();
    console.log("Form submitted");

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
            const data = await response.json();
            localStorage.setItem("username", username);
            localStorage.setItem("token", data.token);

            // check if email is provided and save it
            if (data.email) {
                localStorage.setItem("email", data.email);
            }

            // redirect to index.html
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

// function to display error message
function displayError(message) {
    let errorElement = document.querySelector(".error-message");
    if (!errorElement) {
        errorElement = document.createElement("div");
        errorElement.className = "error-message";
        document.querySelector(".login-container").appendChild(errorElement);
    }
    errorElement.textContent = message;
}