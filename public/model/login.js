document.addEventListener("DOMContentLoaded", () => {
  const loginBox = document.getElementById("loginBox");
  const signupBox = document.getElementById("signupBox");
  const showSignup = document.getElementById("showSignup");
  const showLogin = document.getElementById("showLogin");

  if (!loginBox || !signupBox || !showSignup || !showLogin) {
    return;
  }

  const openSignup = () => {
    loginBox.classList.add("is-hidden");
    signupBox.classList.remove("is-hidden");
  };

  const openLogin = () => {
    signupBox.classList.add("is-hidden");
    loginBox.classList.remove("is-hidden");
  };

  showSignup.addEventListener("click", (event) => {
    event.preventDefault();
    openSignup();
  });

  showLogin.addEventListener("click", (event) => {
    event.preventDefault();
    openLogin();
  });

  openLogin();
});
