<script lang="ts">
    import { addToast } from "../stores/toast-wrapper";
    import type { Community } from "../types/community";

    let {
        searchQuery = $bindable(),
        selectedCategories = $bindable<string[]>([]),
        availableCommunities = $bindable<Community[]>([]),
        filteredCommunities = $bindable<Community[]>([]),
        categories,
    }: {
        searchQuery?: string;
        selectedCategories?: string[];
        availableCommunities?: Community[];
        filteredCommunities?: Community[];
        categories: string[];
    } = $props();

    function toggleCategory(category: string) {
        if (selectedCategories.includes(category)) {
            selectedCategories = selectedCategories.filter(
                (c) => c !== category
            );
        } else {
            selectedCategories = [...selectedCategories, category];
        }
    }

    async function handleSearch() {
        if (searchQuery?.trim() === "" && selectedCategories.length === 0) {
            filteredCommunities = availableCommunities;
            return;
        }

        filteredCommunities = availableCommunities.filter((community) => {
            const matchesQuery =
                !searchQuery ||
                community.community_name
                    .toLowerCase()
                    .includes(searchQuery.toLowerCase());
            const matchesCategory =
                selectedCategories.length === 0 ||
                selectedCategories.some((cat) =>
                    community.categories?.includes(cat)
                );
            return matchesQuery && matchesCategory;
        });
    }

    $effect(() => {
        $inspect(
            "SearchFilter",
            {
                searchQuery,
                selectedCategories,
                availableCommunities,
                filteredCommunities,
            }
        );
    })
</script>

<div class="search-filter">
    <div class="search-bar">
        <input
            type="text"
            placeholder="Search communities..."
            bind:value={searchQuery}
            onchange={() => handleSearch()}
            onkeydown={(e) => e.key === "Enter" && handleSearch()}
        />
    </div>

    <div class="category-filters">
        <span>Filter by categories:</span>
        {#each categories as category}
            <button
                class="category-filter {selectedCategories.includes(category)
                    ? 'active'
                    : ''}"
                onclick={() => toggleCategory(category)}
                onchange={handleSearch}
            >
                {category}
            </button>
        {/each}
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/community.scss";
</style>
