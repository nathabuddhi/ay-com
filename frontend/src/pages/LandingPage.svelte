<script lang="ts">
    import { onMount } from "svelte";
    import ToggleTheme from "../components/ToggleTheme.svelte";
    import Footer from "../components/Footer.svelte";
    import { isLoggedIn } from "../controllers/token-controller";
    import "../styles/app.scss";

    let currentTheme: "light" | "dark" = "dark";
    onMount(async () => {
        if (await isLoggedIn()) {
            window.location.href = "/home";
        }

        const savedTheme =
            (localStorage.getItem("theme") as "light" | "dark") || "light";
        currentTheme = savedTheme;
        document.documentElement.setAttribute("data-theme", savedTheme);
    });
</script>

<div class="landing-container">
    <div class="landing-content">
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
                <h1>Happening now</h1>
                <h2>Join today.</h2>
                <p class="tagline">Connect, share, engage.</p>

                <div class="account-options">
                    <div class="create-account-section">
                        <p class="account-prompt">Don't have an account?</p>
                        <a href="/register" class="create-account-button">
                            Create account
                        </a>
                        <p class="terms">
                            By signing up, you agree to the <a href="/terms"
                                >Terms of Service</a
                            >
                            and
                            <a href="/privacy">Privacy Policy</a>, including
                            <a href="/cookies">Cookie Use</a>.
                        </p>
                    </div>

                    <div class="sign-in-section">
                        <p class="account-prompt">Already have an account?</p>
                        <a href="/login" class="sign-in-button"> Sign in </a>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <Footer />
</div>

<style lang="scss">
    @use "../styles/landing.scss";
</style>
