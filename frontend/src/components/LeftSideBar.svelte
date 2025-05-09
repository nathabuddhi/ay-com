<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import ToggleTheme from "./ToggleTheme.svelte";
    import { BadgeCheck } from "@lucide/svelte";
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

    export let userData: {
        name: string;
        username: string;
        is_verified: boolean;
        user_id: string;
    };

    let showLogoutMenu = false;

    function toggleLogoutMenu(): void {
        showLogoutMenu = !showLogoutMenu;
    }
</script>

<aside class="left-sidebar">
    <div class="sidebar-content">
        <div class="logo">
            <a href="/home">AY</a>
            <ToggleTheme />
        </div>

        <nav class="nav-links">
            <a href="/home" class="nav-link active">
                <span class="icon">
                    <House />
                </span>
                <span class="link-text">Home</span>
            </a>

            <a href="/explore" class="nav-link">
                <span class="icon">
                    <Search />
                </span>
                <span class="link-text">Explore</span>
            </a>

            <a href="/notifications" class="nav-link">
                <span class="icon">
                    <Bell />
                </span>
                <span class="link-text">Notifications</span>
            </a>

            <a href="/messages" class="nav-link">
                <span class="icon">
                    <Mail />
                </span>
                <span class="link-text">Messages</span>
            </a>

            <a href="/bookmarks" class="nav-link">
                <span class="icon">
                    <Bookmark />
                </span>
                <span class="link-text">Bookmarks</span>
            </a>

            <a href="/communities" class="nav-link">
                <span class="icon">
                    <Users />
                </span>
                <span class="link-text">Communities</span>
            </a>

            <a href="/premium" class="nav-link">
                <span class="icon">
                    <Star />
                </span>
                <span class="link-text">Premium</span>
            </a>

            <a href="/profile" class="nav-link">
                <span class="icon">
                    <CircleUserRound />
                </span>
                <span class="link-text">Profile</span>
            </a>

            <a href="/settings" class="nav-link">
                <span class="icon">
                    <Settings />
                </span>
                <span class="link-text">Settings</span>
            </a>
        </nav>

        <button class="post-button">Post</button>
    </div>

    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        class="user-profile"
        on:click={toggleLogoutMenu}
        aria-expanded={showLogoutMenu}
    >
        <div class="profile-image">
            <img
                src={`${import.meta.env.AVATAR_LINK}/${userData.user_id}.png`}
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
                <button on:click={logout}>Log out @{userData.username}</button>
            </div>
        {/if}
    </div>
</aside>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
