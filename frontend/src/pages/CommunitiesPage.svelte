<script lang="ts">
    import { onMount } from "svelte";
    import type { Community } from "../types/community";
    import {
        getAllCommunities,
        getJoinedCommunities,
        getPendingJoinCommunities,
        searchCommunities,
        sendJoinRequest,
        getCommunityCategories,
    } from "../controllers/community-controller";
    import CommunityComponent from "../components/CommunityComponent.svelte";
    import SearchFilter from "../components/SearchFilter.svelte";
    import { addToast } from "../stores/toast-wrapper";
    import Pagination from "../components/Pagination.svelte";

    let items = Array.from({ length: 137 }, (_, i) => i + 1);

    let joinedCommunities = $state<Community[]>([]);
    let pendingCommunities = $state<Community[]>([]);
    let allCommunities = $state<Community[]>([]);
    let filteredCommunities = $state<Community[]>([]);
    let categories = $state<string[]>([]);
    let loading = $state(false);
    let activeTab = $state<"joined" | "pending" | "all">("joined");

    let searchQuery = $state("");
    let selectedCategories = $state<string[]>([]);

    onMount(async () => {
        await loadData();
    });

    async function loadData() {
        loading = true;
        const [joinedRes, pendingRes, availableRes, categoriesRes] =
            await Promise.all([
                getJoinedCommunities(),
                getPendingJoinCommunities(),
                getAllCommunities(),
                getCommunityCategories(),
            ]);

        if (joinedRes.success && joinedRes.payload?.communities) {
            joinedCommunities = joinedRes.payload.communities;
        } else {
            if (!joinedRes.success)
                addToast(
                    "error",
                    joinedRes.message || "Failed to load joined communities",
                    "Error!"
                );
            joinedCommunities = [];
        }

        if (pendingRes.success && pendingRes.payload?.communities) {
            pendingCommunities = pendingRes.payload.communities;
        } else {
            if (!pendingRes.success)
                addToast(
                    "error",
                    pendingRes.message || "Failed to load pending communities",
                    "Error!"
                );
            pendingCommunities = [];
        }

        if (availableRes.success && availableRes.payload?.communities) {
            allCommunities = availableRes.payload.communities;
        } else {
            addToast(
                "error",
                availableRes.message || "Failed to load available communities",
                "Error!"
            );
            allCommunities = [];
        }

        if (categoriesRes.success && categoriesRes.payload?.categories) {
            categories = categoriesRes.payload.categories;
        } else {
            addToast(
                "error",
                categoriesRes.message || "Failed to load categories",
                "Error!"
            );
            categories = [];
        }

        loading = false;
    }
</script>

<div class="communities-page">
    <div class="header">
        <h1>Communities</h1>
        <button
            class="create-btn"
            onclick={() => (window.location.href = "/community/create")}
        >
            Create Community
        </button>
    </div>

    <div class="tabs">
        <button
            class="tab-button {activeTab === 'joined' ? 'active' : ''}"
            onclick={() => (activeTab = "joined")}
        >
            Joined Communities
        </button>
        <button
            class="tab-button {activeTab === 'pending' ? 'active' : ''}"
            onclick={() => (activeTab = "pending")}
        >
            Pending Requests
        </button>
        <button
            class="tab-button {activeTab === 'all' ? 'active' : ''}"
            onclick={() => (activeTab = "all")}
        >
            All Communities
        </button>
    </div>

    {#if activeTab === "all"}
        <SearchFilter
            bind:searchQuery
            bind:selectedCategories
            bind:availableCommunities={allCommunities}
            bind:filteredCommunities
            {categories}
        />
    {:else if activeTab === "joined"}
        <SearchFilter
            bind:searchQuery
            bind:selectedCategories
            bind:availableCommunities={joinedCommunities}
            bind:filteredCommunities
            {categories}
        />{:else if activeTab === "pending"}
        <SearchFilter
            bind:searchQuery
            bind:selectedCategories
            bind:availableCommunities={pendingCommunities}
            bind:filteredCommunities
            {categories}
        />{/if}

    <div class="content">
        {#if loading}
            <div class="loading">Loading communities...</div>
        {:else if searchQuery !== "" || selectedCategories.length > 0}
            <div class="communities-grid">
                {#each filteredCommunities as community}
                    <CommunityComponent {community} />
                {/each}
            </div>
            {#if filteredCommunities.length === 0}
                <p class="no-communities-message">
                    No communities found matching your search.
                </p>
            {/if}
            <Pagination totalItems={filteredCommunities.length ?? 0} />
        {:else if activeTab === "joined"}
            <div class="communities-grid">
                {#each joinedCommunities as community}
                    <CommunityComponent {community} />
                {/each}
            </div>
            {#if filteredCommunities.length === 0}
                <p class="no-communities-message">
                    You have not joined any communities yet.
                </p>
            {/if}
            <Pagination totalItems={joinedCommunities.length ?? 0} />
        {:else if activeTab === "pending"}
            <div class="communities-grid">
                {#each pendingCommunities as community}
                    <CommunityComponent {community} />
                {/each}
            </div>
            {#if pendingCommunities.length === 0}
                <p class="no-communities-message">
                    You do not have any pending join requests yet.
                </p>
            {/if}
            <Pagination totalItems={pendingCommunities.length ?? 0} />
        {:else}
            <div class="communities-grid">
                {#each allCommunities as community}
                    <CommunityComponent {community} />
                {/each}
            </div>
            {#if allCommunities.length === 0}
                <p class="no-communities-message">
                    There are no communities yet.
                </p>
            {/if}
            <Pagination totalItems={allCommunities.length ?? 0} />
        {/if}
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/community.scss";
</style>
