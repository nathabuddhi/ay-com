<script lang="ts">
    import ToggleTheme from "../components/ToggleTheme.svelte";
    import "../styles/app.scss";
    import { MessageCircleQuestion } from "@lucide/svelte";
    import { theme } from "../stores/theme-wrapper";

    function goToHome(): void {
        window.location.href = "/";
    }
</script>

<div class="not-found-container">
    <header>
        <div class="logo-container">
            <img
                src={$theme === "dark" ? "logo-dark.png" : "logo-light.png"}
                alt="AY Logo"
                class="logo-img"
            />
        </div>
        <div class="theme-toggle-container">
            <ToggleTheme />
        </div>
    </header>

    <main>
        <div class="content">
            <h1>404</h1>
            <h2>Page not found</h2>
            <p>The page you're looking for doesn't exist or has been moved.</p>
            <button class="home-button" on:click={goToHome}>
                Return to home
            </button>
        </div>

        <div>
            <MessageCircleQuestion size="300" color="var(--primary-light)" />
        </div>
    </main>
</div>

<style lang="scss">
    @import "../styles/not-found.scss";
</style>
