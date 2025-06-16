<script lang="ts">
    import { ClockIcon, SearchIcon } from "@lucide/svelte";
    import { onMount } from "svelte";
    import ExploreMic from "./ExploreMic.svelte";

    let searchQuery = $state("");
    let isSearchFocused = $state(false);
    let recentSearches: string[] = $state([]);

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
            onfocus={() => (isSearchFocused = true)}
            onblur={() => setTimeout(() => (isSearchFocused = false), 200)}
            onkeydown={handleSearch}
        />
        <ExploreMic bind:query={searchQuery} />
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
                            <ClockIcon />
                        </div>
                        <span>{search}</span>
                    </li>
                {/each}
            </ul>
        </div>
    {/if}
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
