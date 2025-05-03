<script lang="ts">
    import { onMount } from "svelte";
    import ToggleTheme from "../components/ToggleTheme.svelte";
    import Footer from "../components/Footer.svelte";
    import { isLoggedIn, setToken } from "../services/auth";
    import type { UserCredentials } from "../types/user";
    import "../styles/app.scss";
    import type { ApiResponse, StringPayload } from "../types/api";

    let currentTheme: "light" | "dark" = "dark";
    let email: string = "";
    let password: string = "";
    let errorMessage: string | null = null;
    let isLoading: boolean = false;

    onMount(() => {
        if (isLoggedIn()) {
            window.location.href = "/home";
        }

        const savedTheme =
            (localStorage.getItem("theme") as "light" | "dark") || "light";
        currentTheme = savedTheme;
        document.documentElement.setAttribute("data-theme", savedTheme);
    });

    async function handleLogin(event: Event) {
        event.preventDefault();
        isLoading = true;
        errorMessage = null;

        const credentials: UserCredentials = { email, password };

        try {
            const response = await fetch("http://localhost:5000/user/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(credentials),
            });
            console.log("response");
            console.log(response);
            const data: ApiResponse<StringPayload> = await response.json();
            console.log("data");
            console.log(data);
            if (!response.ok || !data.success) {
                throw new Error(data.message || "Invalid email or password");
            }

            setToken(data.payload.value);
            window.location.href = "/home";
        } catch (error) {
            console.log(error);
            errorMessage =
                error instanceof Error
                    ? error.message
                    : "An error occurred during login.";
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="login-container">
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

                {#if errorMessage}
                    <p class="error-message">{errorMessage}</p>
                {/if}

                <form on:submit|preventDefault={handleLogin} class="login-form">
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
                    <p class="account-prompt">Don't have an account?</p>
                    <a href="/register" class="create-account-button">
                        Create account
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
