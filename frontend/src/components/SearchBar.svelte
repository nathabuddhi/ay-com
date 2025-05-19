<script lang="ts">
    import { SearchIcon } from "@lucide/svelte";
    import { onMount } from "svelte";

    let searchQuery = "";
    let isSearchFocused = false;
    let recentSearches: string[] = [];

    onMount(() => {
        const savedSearches = localStorage.getItem("recentSearches");
        if (savedSearches) {
            recentSearches = JSON.parse(savedSearches);
        }
    });

    function handleSearch(event: KeyboardEvent) {
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

    function handleSearchIconClick() {
        if (searchQuery.trim()) {
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
        <div class="search-icon" onclick={handleSearchIconClick}>
            <SearchIcon />
        </div>

        <input
            type="text"
            placeholder="Search"
            bind:value={searchQuery}
            onfocus={() => (isSearchFocused = true)}
            onblur={() => setTimeout(() => (isSearchFocused = false), 200)}
            onkeydown={handleSearch}
        />
    </div>

    {#if isSearchFocused && recentSearches.length > 0}
        <div class="recent-searches">
            <div class="recent-searches-header">
                <h3>Recent searches</h3>
                <button class="clear-button" onclick={clearRecentSearches}
                    >Clear all</button
                >
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
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <circle cx="12" cy="12" r="10"></circle>
                                <polyline points="12 6 12 12 16 14"></polyline>
                            </svg>
                        </div>
                        <span>{search}</span>
                    </li>
                {/each}
            </ul>
        </div>
    {/if}
</div>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
