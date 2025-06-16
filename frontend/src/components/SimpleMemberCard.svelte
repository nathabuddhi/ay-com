<script lang="ts">
    import { AVATAR_IMG } from "../env_var";
    import type { UserProfile } from "../types/user";

    let { user }: { user: UserProfile } = $props();
</script>

<div class="member-item">
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <div
        class="member-info"
        role="button"
        tabindex="0"
        onclick={() => (window.location.href = `/profile/${user.user_id}`)}
    >
        <img
            src={AVATAR_IMG + user.user_id + ".png" || "/placeholder.svg"}
            alt={user.username}
        />
        <span>{user.username}</span>
    </div>
</div>

<style lang="scss">
    .member-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem;
        border-bottom: 1px solid #f8f9fa;

        &:last-child {
            border-bottom: none;
        }

        .member-info {
            display: flex;
            align-items: center;
            gap: 1rem;
            cursor: pointer;
            flex: 1;

            &:hover {
                color: var(--primary);
            }

            img {
                width: 40px;
                height: 40px;
                border-radius: 50%;
                object-fit: cover;
            }
        }
    }
</style>
