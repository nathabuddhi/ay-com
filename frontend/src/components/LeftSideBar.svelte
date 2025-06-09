<script lang="ts">
    import ToggleTheme from "./ToggleTheme.svelte";
    import { BadgeCheck, Menu, ShieldUser } from "@lucide/svelte";
    import { logout } from "../controllers/token-controller";
    import {
        Mail,
        House,
        Search,
        Bell,
        Bookmark,
        Users,
        Star,
        CircleUserRound,
        Settings,
    } from "@lucide/svelte";
    import { onMount } from "svelte";

    let userData = $state({
        name: "",
        username: "",
        is_verified: false,
        user_id: "1",
    });

    onMount(() => {
        userData.name = localStorage.getItem("name") || "User";
        userData.username = localStorage.getItem("username") || "user";
        userData.is_verified = localStorage.getItem("is_verified") === "true";
        userData.user_id = localStorage.getItem("user_id") || "1";
    });

    let showLogoutMenu = $state(false);
    let mobileNavOpen = $state(true);

    function toggleLogoutMenu(): void {
        showLogoutMenu = !showLogoutMenu;
    }

    function toggleMobileNav(): void {
        mobileNavOpen = !mobileNavOpen;
    }

    let { isOpen = $bindable() } = $props();

    const path = $state(window.location.pathname);
</script>

<aside class="left-sidebar {mobileNavOpen ? 'open' : ''}">
    <div class="sidebar-content">
        <div class="logo">
            <button
                class="hamburger"
                onclick={toggleMobileNav}
                aria-label="Toggle navigation"
            >
                <Menu />
            </button>

            <a href="/home">AY</a>
        </div>

        <nav class="nav-links">
            <div class="nav-link">
                <ToggleTheme />
            </div>
            <a
                href="/home"
                class={"nav-link" + (path === "/home" ? " active" : "")}
            >
                <span class="icon">
                    <House />
                </span>
                <span class="link-text">Home</span>
            </a>

            <a
                href="/explore"
                class={"nav-link" + (path === "/explore" ? " active" : "")}
            >
                <span class="icon">
                    <Search />
                </span>
                <span class="link-text">Explore</span>
            </a>

            <a
                href="/notifications"
                class={"nav-link" +
                    (path === "/notifications" ? " active" : "")}
            >
                <span class="icon">
                    <Bell />
                </span>
                <span class="link-text">Notifications</span>
            </a>

            <a
                href="/messages"
                class={"nav-link" + (path === "/message" ? " active" : "")}
            >
                <span class="icon">
                    <Mail />
                </span>
                <span class="link-text">Messages</span>
            </a>

            <a
                href="/bookmarks"
                class={"nav-link" + (path === "/bookmark" ? " active" : "")}
            >
                <span class="icon">
                    <Bookmark />
                </span>
                <span class="link-text">Bookmarks</span>
            </a>

            <a
                href="/communities"
                class={"nav-link" + (path === "/community" ? " active" : "")}
            >
                <span class="icon">
                    <Users />
                </span>
                <span class="link-text">Communities</span>
            </a>

            {#if localStorage.getItem("is_admin") !== "true"}
                <a
                    href="/premium"
                    class={"nav-link" + (path === "/premium" ? " active" : "")}
                >
                    <span class="icon">
                        <Star />
                    </span>
                    <span class="link-text">Premium</span>
                </a>
            {/if}

            <a
                href="/profile"
                class={"nav-link" + (path === "/profile" ? " active" : "")}
            >
                <span class="icon">
                    <CircleUserRound />
                </span>
                <span class="link-text">Profile</span>
            </a>

            {#if localStorage.getItem("is_admin") === "true"}
                <a
                    href="/admin"
                    class={"nav-link" + (path === "/admin" ? " active" : "")}
                >
                    <span class="icon">
                        <ShieldUser />
                    </span>
                    <span class="link-text">Admin</span>
                </a>
            {/if}

            <a href="/settings" class="nav-link">
                <span class="icon">
                    <Settings />
                </span>
                <span class="link-text">Settings</span>
            </a>
        </nav>

        <!-- svelte-ignore element_invalid_self_closing_tag -->
        <button
            aria-label="createpost"
            class="post-button"
            onclick={() => {
                isOpen = true;
            }}
        />
    </div>

    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        class="user-profile"
        onclick={toggleLogoutMenu}
        aria-expanded={showLogoutMenu}
    >
        <div class="profile-image">
            <img
                src={`${import.meta.env.VITE_AVATAR_LINK}${userData.user_id}.png`}
                alt="Profile"
            />
        </div>
        <div class="profile-info">
            <div class="profile-name">
                {userData.name}
                {#if userData.is_verified}
                    <BadgeCheck />
                {/if}
            </div>
            <div class="profile-username">@{userData.username}</div>
        </div>

        {#if showLogoutMenu}
            <div class="logout-menu">
                <button onclick={logout}>Log out @{userData.username}</button>
            </div>
        {/if}
    </div>
</aside>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
