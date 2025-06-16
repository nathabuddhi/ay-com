<script lang="ts">
    import { followUser } from "../controllers/user-controller";
    import { AVATAR_IMG } from "../env_var";
    import { addToast } from "../stores/toast-wrapper";
    import type { UserProfile } from "../types/user";

    let { user }: { user: UserProfile } = $props();

    async function handleFollowClick() {
        const response = await followUser(user.user_id);

        if (response.success) {
            addToast("success", "Followed successfully!", "Success!");
            setTimeout(() => {
                window.location.reload();
            }, 500);
        } else {
            addToast("error", response.message, "Error!");
        }
    }
</script>

<div class="member-item">
    <div class="member-info" role="button" tabindex="0">
        <img
            src={AVATAR_IMG + user.user_id + ".png" || "/placeholder.svg"}
            alt={user.username}
        />
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <span
            onclick={() => (window.location.href = `/profile/${user.username}`)}
        >
            {user.username}
        </span>
        <button class="follow-btn" onclick={handleFollowClick}
            >Follow User</button
        >
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

        .follow-btn {
            background-color: var(--primary);
            color: white;
            border: none;
            padding: 0.5rem 1rem;
            border-radius: 4px;
            cursor: pointer;
            transition: background-color 0.3s;

            &:hover {
                background-color: var(--primary-hover);
            }
        }
    }
</style>
