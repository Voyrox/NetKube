document.addEventListener("DOMContentLoaded", () => {
  const loginForm = document.getElementById("loginForm");
  const loginButton = document.getElementById("loginButton");
  const loginStatus = document.getElementById("loginStatus");
  if (!loginForm || !loginButton || !loginStatus) return;

  const setStatus = (message, isError) => {
    loginStatus.textContent = message;
    loginStatus.classList.toggle("is-visible", Boolean(message));
    loginStatus.classList.toggle("is-error", Boolean(message && isError));
    loginStatus.classList.toggle("is-success", Boolean(message && !isError));
  };

  loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    setStatus("", false);

    const formData = new FormData(loginForm);
    const email = String(formData.get("email") || "").trim();
    const password = String(formData.get("password") || "");

    loginButton.disabled = true;
    loginButton.textContent = "Logging in...";
    try {
      const response = await fetch("/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });
      const payload = await response.json().catch(() => ({}));
      if (!response.ok) {
        throw new Error(payload.error || "Login failed.");
      }

      setStatus("Login successful. Redirecting...", false);
      window.location.href = payload.redirect || "/";
    } catch (error) {
      setStatus(error.message || "Login failed.", true);
    } finally {
      loginButton.disabled = false;
      loginButton.textContent = "Log In";
    }
  });
});
