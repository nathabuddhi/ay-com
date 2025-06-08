<script lang="ts">
    import { BadgeCheck } from "@lucide/svelte";
    import type { UserProfile } from "../types/user";
    import { onMount } from "svelte";
    import {
        followUser,
        getFollowRecommendations,
        getUserById,
    } from "../controllers/user-controller";
    import { AVATAR_IMG } from "../env_var";
    import { addToast } from "../stores/toast-wrapper";
    import { navigate } from "svelte-routing";

    let userids = $state<string[]>([]);
    let users = $state<UserProfile[]>([]);

    onMount(async () => {
        const response = await getFollowRecommendations();

        if (response.success && response.payload) {
            userids = response.payload.user_ids;

            for (const userId of userids) {
                const userResponse = await getUserById(userId);
                if (userResponse) {
                    users.push(userResponse);
                }
            }
        } else {
            userids = [];
            users = [];
        }
    });

    async function handleFollow(user: UserProfile) {
        const response = await followUser(user.user_id);

        if (response.success) {
            addToast("success", "Followed successfully!", "Success!");
            setTimeout(() => {
                navigate(`/profile/${user.username}`);
            }, 500);
        } else {
            addToast("error", response.message, "Error!");
        }
    }
</script>

<div class="who-to-follow-section">
    <h2 class="section-title">Who to follow</h2>

    <div class="who-to-follow-list">
        {#each users as user}
            <div class="who-to-follow-item">
                <div class="user-avatar">
                    <img
                        src={`${AVATAR_IMG}/${user.user_id}.png`}
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

                <button
                    class="follow-button"
                    onclick={() => {
                        handleFollow(user);
                    }}>Follow</button
                >
            </div>
        {/each}
    </div>

    <a href="/connect" class="show-more">Show more</a>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
