<script lang="ts">
    import { EllipsisVertical } from "@lucide/svelte";
    import type { Hashtag } from "../types/thread";
    import { onMount } from "svelte";
    import { getTrendingTags } from "../controllers/thread-controller";

    let hashtags: Hashtag[] = $state<Hashtag[]>([]);

    function formatNumber(num: number): string {
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + "M";
        } else if (num >= 1000) {
            return (num / 1000).toFixed(1) + "K";
        } else {
            return num.toString();
        }
    }

    onMount(async () => {
        const response = await getTrendingTags();

        if (response.success && response.payload) {
            hashtags = response.payload.hashtags;
        } else {
            hashtags = [];
        }
    });
</script>

<div class="trending-section">
    <h2 class="section-title">What's happening</h2>

    <div class="trending-list">
        {#each hashtags as hashtag}
            <a href={`/explore?q=%23${hashtag.hashtag}`} class="trending-item">
                <div class="trending-tag">
                    <span class="trending-label">Trending</span>
                    <h3 class="trending-hashtag">#{hashtag.hashtag}</h3>
                    <span class="trending-count"
                        >{formatNumber(hashtag.thread_count)} posts</span
                    >
                </div>
            </a>
        {/each}
        {#if hashtags.length === 0}
            <div class="no-trending">
                <p>No trending topics at the moment.</p>
            </div>
        {/if}
    </div>

    <a href="/explore" class="show-more">Show more</a>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
