<script lang="ts">
    import type { Thread } from "../types/thread";
    import Post from "./Post.svelte";
    // import CreatePost from "./CreatePost.svelte";

    const forYouPosts = $state<Thread[]>([]);
    const followingPosts = $state<Thread[]>([]);

    const tabs = ["For you", "Following"];
    let activeTab = $state<string>("For you");

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
                    onclick={() => setActiveTab(tab)}
                >
                    {tab}
                </button>
            {/each}
        </div>
    </div>

    <div class="posts-container">
        {#if activeTab === "For you"}
            {#each forYouPosts as post}
                <Post {post} />
            {/each}
            {#if forYouPosts.length === 0}
                <p class="no-posts-message">No posts yet.</p>
            {:else}
                <p class="end-posts-message">You have reached the end.</p>
            {/if}
        {:else}
            {#each followingPosts as post}
                <Post {post} />
            {/each}
            {#if followingPosts.length === 0}
                <p class="no-posts-message">No posts yet.</p>
            {:else}
                <p class="end-posts-message">You have reached the end.</p>
            {/if}
        {/if}
    </div>
</div>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
