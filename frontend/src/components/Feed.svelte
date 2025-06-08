<script lang="ts">
    import { onMount } from "svelte";
    import type { Thread } from "../types/thread";
    import Post from "./Post.svelte";
    import {
        getFollowingThreads,
        getForYouThreads,
    } from "../controllers/thread-controller";
    import { addToast } from "../stores/toast-wrapper";
    import CreatePostForm from "./CreatePostForm.svelte";
    // import CreatePost from "./CreatePost.svelte";

    let forYouPosts = $state<Thread[]>([]);
    let followingPosts = $state<Thread[]>([]);

    const tabs = ["For you", "Following"];
    let activeTab = $state<string>("For you");

    function setActiveTab(tab: string) {
        activeTab = tab;
    }

    onMount(async () => {
        const response = await getForYouThreads();

        if (response.success) {
            forYouPosts = response.payload?.threads || [];
        } else {
            addToast(
                "error",
                "Failed to fetch general feed: " + response.message,
                "Error!"
            );
        }

        const response2 = await getFollowingThreads();

        if (response2.success) {
            followingPosts = response2.payload?.threads || [];
        } else {
            addToast(
                "error",
                "Failed to fetch following feed: " + response.message,
                "Error!"
            );
        }
    });
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
    <CreatePostForm isOpen={true} mode={"post"} />
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

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
