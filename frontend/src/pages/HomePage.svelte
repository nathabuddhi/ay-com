<script lang="ts">
    import { onMount } from "svelte";
    import LeftSideBar from "../components/LeftSideBar.svelte";
    import RightBar from "../components/RightBar.svelte";
    import Feed from "../components/Feed.svelte";
    import { isLoggedIn } from "../controllers/token-controller";
    import "../styles/app.scss";
    import ProfilePage from "./ProfilePage.svelte";
    import ToastContainer from "../components/ToastContainer.svelte";
    import NotificationPage from "./NotificationPage.svelte";
    import CreatePost from "../components/CreatePost.svelte";
    import BookmarkPage from "./BookmarkPage.svelte";
    import ThreadDetailPage from "./ThreadDetailPage.svelte";

    onMount(async () => {
        if (!(await isLoggedIn())) {
            window.location.href = "/login";
            return;
        }
    });

    const currpage = window.location.pathname;
    let isNewPostOpen = false;
</script>

<div class="home-container">
    <ToastContainer />
    <CreatePost bind:isOpen={isNewPostOpen} />
    <LeftSideBar bind:isNewPostOpen />
    <main class="main-content">
        {#if currpage.includes("home")}
            <Feed />
        {:else if currpage.includes("profile")}
            <ProfilePage />
        {:else if currpage.includes("notifications")}
            <NotificationPage />
        {:else if currpage.includes("bookmarks")}
            <BookmarkPage />
        {:else if currpage.includes("thread")}
            <ThreadDetailPage />
        {/if}
    </main>
    <RightBar />
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
