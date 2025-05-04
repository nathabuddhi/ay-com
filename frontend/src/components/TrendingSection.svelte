<script lang="ts">
    export let hashtags: {
        tag: string;
        count: number;
    }[];
    
    // Format numbers (e.g., 1.2K instead of 1200)
    function formatNumber(num: number): string {
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + "M";
        } else if (num >= 1000) {
            return (num / 1000).toFixed(1) + "K";
        } else {
            return num.toString();
        }
    }
</script>

<div class="trending-section">
    <h2 class="section-title">What's happening</h2>
    
    <div class="trending-list">
        {#each hashtags as hashtag}
            <a href={`/explore?q=%23${hashtag.tag}`} class="trending-item">
                <div class="trending-tag">
                    <span class="trending-label">Trending</span>
                    <h3 class="trending-hashtag">#{hashtag.tag}</h3>
                    <span class="trending-count">{formatNumber(hashtag.count)} posts</span>
                </div>
                
                <!-- svelte-ignore a11y_consider_explicit_label -->
                <button class="more-options">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="12" cy="12" r="1"></circle>
                        <circle cx="12" cy="5" r="1"></circle>
                        <circle cx="12" cy="19" r="1"></circle>
                    </svg>
                </button>
            </a>
        {/each}
    </div>
    
    <a href="/explore" class="show-more">Show more</a>
</div>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>