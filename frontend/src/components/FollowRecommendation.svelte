<script lang="ts">
    import { BadgeCheck } from "@lucide/svelte";
    import type { UserProfile } from "../types/user";

    let users = $state<UserProfile[]>([]);
</script>

<div class="who-to-follow-section">
    <h2 class="section-title">Who to follow</h2>

    <div class="who-to-follow-list">
        {#each users as user}
            <div class="who-to-follow-item">
                <div class="user-avatar">
                    <img
                        src={`http://localhost:5000/images/profile/${user.user_id}`}
                        alt={user.name}
                    />
                </div>

                <div class="user-info">
                    <div class="user-name">
                        {user.name}
                        {#if user.is_verified}
                            <BadgeCheck />
                        {/if}
                    </div>
                    <div class="user-username">@{user.username}</div>
                </div>

                <button class="follow-button">Follow</button>
            </div>
        {/each}
    </div>

    <a href="/connect" class="show-more">Show more</a>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
