<script lang="ts">
    import SearchBar from "./SearchBar.svelte";
    import TrendingSection from "./TrendingSection.svelte";
    import FollowRecommendations from "./FollowRecommendation.svelte";
</script>

<aside class="right-sidebar">
    <div class="sidebar-content">
        <SearchBar />

        {#if localStorage.getItem("is_verified") === null || localStorage.getItem("is_verified") === "false"}
            <div class="premium-promo">
                <h2 class="promo-title">Subscribe to Premium</h2>
                <p class="promo-description">
                    Subscribe to unlock new features and if eligible, receive a
                    share of revenue.
                </p>
                <a href="/premium" class="subscribe-button">Subscribe</a>
            </div>
        {/if}

        <TrendingSection />

        <FollowRecommendations />
    </div>
</aside>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
