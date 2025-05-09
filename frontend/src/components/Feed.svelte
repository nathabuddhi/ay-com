<script lang="ts">
    import type { Thread } from "../types/thread";
    import Post from "./Post.svelte";
    // import CreatePost from "./CreatePost.svelte";

    const posts: Thread[] = [];

    const tabs = ["For you", "Following"];
    let activeTab = "For you";

    function setActiveTab(tab: string) {
        activeTab = tab;
    }
</script>

<div class="feed-container">
    <div class="feed-header">
        <h1>Home</h1>
        <div class="feed-tabs">
            {#each tabs as tab}
                <button
                    class="tab-button {activeTab === tab ? 'active' : ''}"
                    on:click={() => setActiveTab(tab)}
                >
                    {tab}
                </button>
            {/each}
        </div>
    </div>

    <div class="posts-container">
        {#each posts as post}
            <Post {post} />
        {/each}
        {#if posts.length === 0}
            <p class="no-posts-message">No posts yet.</p>
        {:else}
            <p class="end-posts-message">You have reached the end.</p>
        {/if}
    </div>
</div>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
