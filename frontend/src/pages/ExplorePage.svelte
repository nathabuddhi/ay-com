<script lang="ts">
    import { onMount } from "svelte";
    import ExploreSearchBar from "../components/ExploreSearchBar.svelte";

    const tabs = ["Top", "Latest", "People", "Media", "Communities"];
    let activeTab = $state("Top");

    function setActiveTab(tab: string) {
        activeTab = tab;
    }

    onMount(() => {});
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
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/explore.scss";
</style>
