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
    import CreatePostDialog from "../components/CreatePostDialog.svelte";
    import BookmarkPage from "./BookmarkPage.svelte";
    import ThreadDetailPage from "./ThreadDetailPage.svelte";
    import PremiumPage from "./PremiumPage.svelte";
    import AdminPage from "./AdminPage.svelte";
    import CreateCommunityPage from "./CreateCommunityPage.svelte";
    import CommunitiesPage from "./CommunitiesPage.svelte";
    import CommunityDetailPage from "./CommunityDetailPage.svelte";
    import ExplorePage from "./ExplorePage.svelte";

    onMount(async () => {
        if (!(await isLoggedIn())) {
            window.location.href = "/login";
            return;
        }
    });

    const currpage = window.location.pathname;

    let isOpen = $state(false);
</script>

<div class="home-container">
    <ToastContainer />
    <CreatePostDialog bind:isOpen />
    <LeftSideBar bind:isOpen />
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
        {:else if currpage.includes("premium")}
            <PremiumPage />
        {:else if currpage.includes("admin")}
            <AdminPage />
        {:else if currpage.includes("community/create")}
            <CreateCommunityPage />
        {:else if currpage.includes("community/")}
            <CommunityDetailPage />
        {:else if currpage.includes("communities")}
            <CommunitiesPage />
        {:else if currpage.includes("explore")}
            <ExplorePage />
        {/if}
    </main>
    <RightBar />
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
