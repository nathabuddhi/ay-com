<script lang="ts">
    import { ClockIcon, SearchIcon } from "@lucide/svelte";
    import ExploreMic from "./ExploreMic.svelte";
    import type { UserProfile } from "../types/user";
    import { searchUsers } from "../controllers/user-controller";
    import { AVATAR_IMG } from "../env_var";
    import { onMount } from "svelte";

    let searchQuery = $state("");
    let isSearchFocused = $state(false);
    let userRecommendations = $state<UserProfile[]>([]);

    let debounceTimeout: ReturnType<typeof setTimeout>;
    let recentSearches: string[] = $state([]);

    onMount(async () => {
        const savedSearches = localStorage.getItem("recentSearches");
        if (savedSearches) {
            recentSearches = JSON.parse(savedSearches);
        }

        searchQuery =
            new URLSearchParams(window.location.search).get("q") || "";
        loadTabData();
    });

    async function loadTabData() {}

    function searchDebouncer() {
        clearTimeout(debounceTimeout);
        debounceTimeout = setTimeout(() => {
            handleSearch(new KeyboardEvent("keydown", { key: "Enter" }));
        }, 300);
    }

    async function handleSearch(event: KeyboardEvent) {
        clearTimeout(debounceTimeout);

        debounceTimeout = setTimeout(async () => {
            const response = await searchUsers(searchQuery);
            if (response.success && response.payload) {
                userRecommendations = response.payload.users;
            } else {
                userRecommendations = [];
            }
        }, 300);

        if (event.key === "Enter" && searchQuery.trim()) {
            if (!recentSearches.includes(searchQuery)) {
                recentSearches = [searchQuery, ...recentSearches.slice(0, 2)];
                localStorage.setItem(
                    "recentSearches",
                    JSON.stringify(recentSearches)
                );
            }
            window.location.href = `/explore?q=${encodeURIComponent(searchQuery)}`;
        }
    }

    function clearRecentSearches() {
        recentSearches = [];
        localStorage.removeItem("recentSearches");
    }

    function selectRecentSearch(search: string) {
        searchQuery = search;
        window.location.href = `/explore?q=${encodeURIComponent(search)}`;
    }
</script>

<div class="search-container">
    <div class="search-input-container">
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div class="search-icon">
            <SearchIcon />
        </div>

        <input
            type="text"
            placeholder="Search"
            bind:value={searchQuery}
            onfocus={() => {
                isSearchFocused = true;
                handleSearch(new KeyboardEvent("keydown", { key: "Space" }));
            }}
            onblur={() => setTimeout(() => (isSearchFocused = false), 200)}
            onkeydown={handleSearch}
            onchange={searchDebouncer}
        />
        <ExploreMic bind:query={searchQuery} />
    </div>

    {#if searchQuery !== "" && isSearchFocused}
        <div class="dropdown-list">
            <ul class="users-list">
                <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
                {#each userRecommendations as user}
                    <!-- svelte-ignore a11y_click_events_have_key_events -->
                    <li
                        class="user-card"
                        onclick={() =>
                            (window.location.href = `/profile/${user.username}`)}
                    >
                        <img
                            src={AVATAR_IMG + user.user_id + ".png"}
                            alt={user.username}
                            class="user-image"
                        />
                        <span class="username">@{user.username}</span>
                    </li>
                {/each}
            </ul>

            {#if recentSearches.length > 0}
                <div class="recent-searches">
                    <div class="recent-searches-header">
                        <h3>Recent searches</h3>
                        <button
                            class="clear-button"
                            onclick={clearRecentSearches}
                        >
                            Clear all
                        </button>
                    </div>

                    <ul class="recent-search-list">
                        <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
                        {#each recentSearches as search}
                            <!-- svelte-ignore a11y_click_events_have_key_events -->
                            <li
                                class="recent-search-item"
                                onclick={() => selectRecentSearch(search)}
                            >
                                <div class="search-history-icon">
                                    <ClockIcon />
                                </div>
                                <span>{search}</span>
                            </li>
                        {/each}
                    </ul>
                </div>
            {/if}
        </div>
    {/if}
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/explore.scss";
</style>
