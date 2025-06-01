<script lang="ts">
    import { onMount } from "svelte";
    import type { Thread } from "../types/thread";
    import Post from "../components/Post.svelte";
    import { getBookmarkedThreads } from "../controllers/thread-controller";
    import { addToast } from "../stores/toast-wrapper";

    let bookmarkedPosts: Thread[] = $state<Thread[]>([]);
    let allBookmarkedPosts: Thread[] = $state<Thread[]>([]);
    let searchQuery = $state<string>("");

    onMount(async () => {
        const response = await getBookmarkedThreads();

        if (response.success) {
            allBookmarkedPosts = response.payload?.threads || [];
            bookmarkedPosts = [...allBookmarkedPosts];
        } else {
            addToast(
                "error",
                "Failed to fetch bookmarked threads: " + response.message,
                "Error!"
            );
        }
    });

    function handleSearch() {
        if (searchQuery.trim() === "") {
            bookmarkedPosts = [...allBookmarkedPosts];
        } else {
            bookmarkedPosts = allBookmarkedPosts.filter((post) =>
                post.content.toLowerCase().includes(searchQuery.toLowerCase())
            );
        }
    }
</script>

<div class="feed-container">
    <div class="feed-header">
        <h1>Bookmarks</h1>
    </div>
    <div class="search-bar">
        <input
            type="text"
            placeholder="Search Bookmarks"
            bind:value={searchQuery}
            oninput={handleSearch}
        />
    </div>
    <div class="posts-container">
        {#each bookmarkedPosts as post}
            <Post {post} />
        {/each}
        {#if bookmarkedPosts.length === 0}
            <p class="no-posts-message">No bookmarks yet.</p>
        {:else}
            <p class="end-posts-message">You have reached the end.</p>
        {/if}
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
    @use "../styles/bookmark.scss";
</style>
