<script lang="ts">
    import { Ban, UserCheck } from "@lucide/svelte";
    import { AVATAR_IMG } from "../../env_var";
    import type { AdminUserProfile } from "../../types/admin";
    import { addToast } from "../../stores/toast-wrapper";
    import { toggleUserBan } from "../../controllers/admin-controller";

    let { user }: { user: AdminUserProfile } = $props();

    async function handleBanUser(userId: string) {
        try {
            const response = await toggleUserBan(userId);
            if (response.success) {
                addToast("success", "User banned successfully");
                user.is_banned = true;
            } else {
                addToast("error", response.message || "Failed to ban user");
            }
        } catch (error) {
            addToast("error", "Failed to ban user");
        }
    }

    async function handleUnbanUser(userId: string) {
        try {
            const response = await toggleUserBan(userId);
            if (response.success) {
                addToast("success", "User unbanned successfully");
                user.is_banned = false;
            } else {
                addToast("error", response.message || "Failed to unban user");
            }
        } catch (error) {
            addToast("error", "Failed to unban user");
        }
    }
</script>

<div class="user-item">
    <div class="user-info">
        <img
            src={AVATAR_IMG + "/" + user.user_id + ".png"}
            alt=""
            class="user-avatar"
        />
        <div class="user-details">
            <div class="user-name">
                <a href="/profile/{user.username}" target="_blank">
                    @{user.username}
                </a>
                <span class="real-name">
                    ({user.name})
                </span>
            </div>
            <div class="user-status">
                {#if user.is_verified}
                    <span class="badge verified"> Verified</span>
                {/if}
                {#if user.is_banned}
                    <span class="badge banned"> Banned </span>
                {/if}
            </div>
        </div>
    </div>
    <div class="user-actions">
        {#if user.is_banned}
            <button
                class="action-btn unban-btn"
                onclick={() => handleUnbanUser(user.user_id)}
            >
                <UserCheck size={16} />
                Unban
            </button>
        {:else}
            <button
                class="action-btn ban-btn"
                onclick={() => handleBanUser(user.user_id)}
            >
                <Ban size={16} />
                Ban
            </button>
        {/if}
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../../styles/admin.scss";
</style>
