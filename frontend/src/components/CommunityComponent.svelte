<script lang="ts">
    import type { Community } from "../types/community";
    import { addToast } from "../stores/toast-wrapper";
    import { sendJoinRequest } from "../controllers/community-controller";
    import { C_ICON_IMG } from "../env_var";

    export let community: Community;

    function handleRedirect() {
        window.location.href = `/community/${community.community_id}`;
    }

    async function handleJoinRequest(e: MouseEvent) {
        e.stopPropagation();
        const response = await sendJoinRequest(community.community_id);
        if (response.success) {
            community.is_pending = true;
            addToast("success", "Sent join request succesfully!");
        } else {
            addToast("error", response.message);
        }
    }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<article class="community-card">
    <div class="logo">
        <img
            src={`${C_ICON_IMG}${community.icon_image}`}
            alt={community.community_name}
        />
    </div>

    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <div class="content" onclick={handleRedirect}>
        <div class="header">
            <div class="post-user-info">
                <span class="post-user-name">{community.community_name}</span>

                <div class="post-time">
                    {#each community.categories as cat}
                        <span class="community-category">{cat}</span>
                    {/each}
                </div>
            </div>

            <button
                class="join-button"
                onclick={handleJoinRequest}
                disabled={(community.role !== "" &&
                    community.role !== "guest") ||
                    community.is_pending}
            >
                {#if community.is_pending}
                    Pending
                {:else if community.role === "member" || community.role === "moderator" || community.role === "owner"}
                    Joined
                {:else}
                    Join
                {/if}
            </button>
        </div>

        <div class="description">{community.description}</div>
    </div>
</article>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/community.scss";
</style>
