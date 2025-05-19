<script lang="ts">
    import { EllipsisVertical } from "@lucide/svelte";

    export let hashtags: {
        tag: string;
        count: number;
    }[];

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
                    <span class="trending-count"
                        >{formatNumber(hashtag.count)} posts</span
                    >
                </div>

                <!-- svelte-ignore a11y_consider_explicit_label -->
                <button class="more-options">
                    <EllipsisVertical />
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
