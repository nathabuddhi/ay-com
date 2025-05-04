<script lang="ts">
    import SearchBar from "./SearchBar.svelte";
    import TrendingSection from "./TrendingSection.svelte";
    import WhoToFollowSection from "./FollowRecommendation.svelte";

    const trendingHashtags = [
        { tag: "gaming", count: 52000 },
        { tag: "valorant", count: 43200 },
        { tag: "minecraft", count: 38500 },
        { tag: "teddyseekinspirasi", count: 25600 },
        { tag: "cr7", count: 18900 },
    ];

    const whoToFollow = [
        {
            id: "1",
            name: "Brainwits",
            username: "brainwits",
            is_verified: true,
            followers: 1200000,
        },
        {
            id: "2",
            name: "Peach",
            username: "peach",
            is_verified: true,
            followers: 980000,
        },
        {
            id: "3",
            name: "Anthony Pompliano",
            username: "APompliano",
            is_verified: true,
            followers: 850000,
        },
    ];
</script>

<aside class="right-sidebar">
    <div class="sidebar-content">
        <SearchBar />

        
        {#if localStorage.getItem('is_verified') === null || localStorage.getItem('is_verified') === 'false'}
            <div class="premium-promo">
            <h2 class="promo-title">Subscribe to Premium</h2>
            <p class="promo-description">
                Subscribe to unlock new features and if eligible, receive a
                share of revenue.
            </p>
            <a href="/premium" class="subscribe-button">Subscribe</a>
            </div>
        {/if}

        <TrendingSection hashtags={trendingHashtags} />

        <WhoToFollowSection users={whoToFollow} />
    </div>
</aside>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
