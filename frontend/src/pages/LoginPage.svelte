<script lang="ts">
    import { onMount } from "svelte";
    import ToggleTheme from "../components/ToggleTheme.svelte";
    import Footer from "../components/Footer.svelte";
    import { isLoggedIn, setToken } from "../controllers/token-controller";
    import "../styles/app.scss";
    import { login } from "../controllers/user-controller";
    import { addToast } from "../stores/toast-wrapper";
    import ToastContainer from "../components/ToastContainer.svelte";

    let currentTheme: "light" | "dark" = "dark";
    let email: string = "";
    let password: string = "";
    let isLoading: boolean = false;

    onMount(async () => {
        if (await isLoggedIn()) {
            addToast(
                "success",
                "User cookie validated! Redirecting...",
                "Success."
            );
            setTimeout(() => {
                window.location.href = "/home";
            }, 2500);
        }

        const savedTheme =
            (localStorage.getItem("theme") as "light" | "dark") || "light";
        currentTheme = savedTheme;
        document.documentElement.setAttribute("data-theme", savedTheme);
    });

    async function onLoginClick(event: Event) {
        event.preventDefault();
        isLoading = true;

        let response = await login(email, password);
        if (response.success) {
            addToast("success", "Login successful! Redirecting...", "Success.");
            setTimeout(() => {
                window.location.href = "/home";
            }, 1500);
        } else {
            addToast("error", response.message, "Failed logging in.");
            isLoading = false;
        }
    }
</script>

<div class="login-container">
    <ToastContainer />
    <div class="login-content">
        <div class="logo-section">
            <img
                src={currentTheme === "dark"
                    ? "logo-dark.png"
                    : "logo-light.png"}
                alt="logo"
                class="logo-img"
            />
        </div>

        <div class="auth-section">
            <div class="theme-toggle-container">
                <ToggleTheme bind:theme={currentTheme} />
            </div>

            <div class="auth-content">
                <h1>Sign in</h1>
                <h2>Access your account</h2>
                <p class="tagline">Connect, share, engage.</p>

                <form
                    on:submit|preventDefault={onLoginClick}
                    class="login-form"
                >
                    <div class="form-group">
                        <label for="email">Email</label>
                        <input
                            type="email"
                            id="email"
                            bind:value={email}
                            required
                            placeholder="Enter your email"
                        />
                    </div>

                    <div class="form-group">
                        <label for="password">Password</label>
                        <input
                            type="password"
                            id="password"
                            bind:value={password}
                            required
                            placeholder="Enter your password"
                        />
                    </div>

                    <button
                        type="submit"
                        class="sign-in-button"
                        disabled={isLoading}
                    >
                        {isLoading ? "Signing in..." : "Sign in"}
                    </button>
                </form>

                <div class="account-options">
                    <a href="/register" class="create-account-button">
                        Don't have an account? Create account
                    </a>
                    <a href="/forgot" class="create-account-button">
                        Forgot your account? Recover it
                    </a>
                </div>
            </div>
        </div>
    </div>

    <Footer />
</div>

<style lang="scss">
    @use "../styles/login.scss";
</style>
