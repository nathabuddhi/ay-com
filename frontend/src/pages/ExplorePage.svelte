<script lang="ts">
    import { onMount } from "svelte";
    import ExploreSearchBar from "../components/ExploreSearchBar.svelte";
    import type { Hashtag, Thread } from "../types/thread";
    import type { UserProfile } from "../types/user";
    import type { Community } from "../types/community";
    import {
        getThreadsByHashtag,
        getTrendingTags,
    } from "../controllers/thread-controller";
    import Post from "../components/Post.svelte";

    let tabs = $state<string[]>([""]);
    let activeTab = $state("Top");
    let query = $state("");

    // TAB DATA
    // Trending
    let hashtags = $state<Hashtag[]>([]);
    let selectedHashtag = $state<string>("");
    let hashtagPost = $state<Thread[]>([]);

    // Top
    let topUsers = $state<UserProfile[]>([]);
    let topPosts = $state<Thread[]>([]);

    // Latest
    let latestPosts = $state<Thread[]>([]);

    // people
    let allUsers = $state<UserProfile[]>([]);

    // media
    let mediaPosts = $state<Thread[]>([]);

    // communities
    let communities = $state<Community[]>([]);

    async function loadAllData() {
        const hashTagResponse = await getTrendingTags();
        if (hashTagResponse.success && hashTagResponse.payload) {
            hashtags = hashTagResponse.payload.hashtags;
        } else {
            hashtags = [];
        }
    }

    async function setSelectedHashtag(hashtag: string) {
        selectedHashtag = hashtag;

        const response = await getThreadsByHashtag(hashtag);

        if (response.success && response.payload) {
            hashtagPost = response.payload.threads;
        } else {
            hashtagPost = [];
        }
    }

    function setActiveTab(tab: string) {
        activeTab = tab;
    }

    onMount(async () => {
        query =
            (await new URLSearchParams(window.location.search).get("q")) || "";
        if (query !== "") {
            tabs = ["Top", "Latest", "People", "Media", "Communities"];
            activeTab = "Top";
        } else {
            tabs = ["Trending"];
            activeTab = "Trending";
        }

        loadAllData();
    });
</script>

<div class="explore-container">
    <div class="explore-header">
        <h1>Explore</h1>
        <ExploreSearchBar />
    </div>
    <div class="tabs">
        {#each tabs as tab}
            <button
                class="tab-button {activeTab === tab ? 'active' : ''}"
                onclick={() => setActiveTab(tab)}
            >
                {tab}
            </button>
        {/each}
    </div>
    <div class="tab-content">
        {#if activeTab == "Trending"}
            <h2>Trending Hashtags</h2>
            <div class="hashtag-list">
                {#each hashtags as hashtag}
                    <button
                        class="hashtag-btn {selectedHashtag === hashtag.hashtag
                            ? 'selected'
                            : ''}"
                        onclick={() => setSelectedHashtag(hashtag.hashtag)}
                    >
                        #{hashtag.hashtag}
                    </button>
                {/each}
            </div>
            {#if selectedHashtag}
                <h3>Posts for #{selectedHashtag}</h3>
                {#if hashtagPost.length > 0}
                    <ul class="hashtag-posts">
                        {#each hashtagPost as post}
                            <Post {post} />
                        {/each}
                    </ul>
                {:else}
                    <p>No posts found for this hashtag.</p>
                {/if}
            {/if}
        {:else if activeTab === "Top"}
            <p>Top content goes here...</p>
        {:else if activeTab === "Latest"}
            <p>Latest content goes here...</p>
        {:else if activeTab === "People"}
            <p>People content goes here...</p>
        {:else if activeTab === "Media"}
            <p>Media content goes here...</p>
        {:else if activeTab === "Communities"}
            <p>Communities content goes here...</p>
        {/if}
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/explore.scss";
</style>
