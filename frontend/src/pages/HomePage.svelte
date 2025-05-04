<script lang="ts">
    import { onMount } from "svelte";
    import LeftSideBar from "../components/LeftSideBar.svelte";
    import RightBar from "../components/RightBar.svelte";
    import Feed from "../components/Feed.svelte";
    import { isLoggedIn } from "../controllers/token-controller";
    import "../styles/app.scss";

    let currentTheme: "light" | "dark" = "dark";
    let userData = {
        name: "",
        username: "",
        is_verified: false,
        user_id: "",
    };

    onMount(async () => {
        if (!(await isLoggedIn())) {
            window.location.href = "/login";
            return;
        }

        // Get user data from localStorage
        userData.name = localStorage.getItem("name") || "User";
        userData.username = localStorage.getItem("username") || "user";
        userData.is_verified = localStorage.getItem("is_verified") === "true";
        userData.user_id = localStorage.getItem("user_id") || "1";

        // Set theme
        const savedTheme =
            (localStorage.getItem("theme") as "light" | "dark") || "light";
        currentTheme = savedTheme;
        document.documentElement.setAttribute("data-theme", savedTheme);
    });

    function handleThemeChange(
        event: CustomEvent<{ theme: "light" | "dark" }>
    ): void {
        currentTheme = event.detail.theme;
        localStorage.setItem("theme", currentTheme);
        document.documentElement.setAttribute("data-theme", currentTheme);
    }
</script>

<div class="home-container">
    <LeftSideBar {userData} {currentTheme} on:themeChange={handleThemeChange} />
    <main class="main-content">
        <Feed />
    </main>
    <RightBar />
</div>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
