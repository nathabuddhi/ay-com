<script lang="ts">
    import { onMount } from "svelte";
    import {
        getBlockedUsers,
        unblockUser,
    } from "../controllers/user-controller";
    import type { BlockedUser } from "../types/user";
    import { addToast } from "../stores/toast-wrapper";

    let blockedUsers = $state<BlockedUser[]>([]);

    async function fetchBlockedUsers() {
        const response = await getBlockedUsers();
        if (response && response.success && response.payload) {
            blockedUsers = response.payload.blocked as unknown as BlockedUser[];
            console.log(response);
        } else {
            addToast(
                "error",
                "Failed to fetch blocked users: " + response?.message,
                "An error occured!"
            );
        }
    }

    onMount(() => {
        fetchBlockedUsers();
    });

    async function handleUnblock(
        user_id: string,
        username: string
    ): Promise<void> {
        addToast("info", `Unblocking user ${username}`, "Unblocking...");

        const response = await unblockUser(user_id);
        if (response && response.success) {
            addToast(
                "success",
                `User ${username} unblocked successfully!`,
                "Success!"
            );
            fetchBlockedUsers();
        } else {
            addToast(
                "error",
                "Failed to unblock user " + username,
                "An error occured!"
            );
        }
    }
</script>

<div class="blocked-users-container">
    <h2 class="title">Blocked Users</h2>
    <div class="users-list">
        {#each blockedUsers as user (user.user_id)}
            <div class="user-item">
                <div class="user-info">
                    <div class="avatar">
                        <img
                            src={`${import.meta.env.AVATAR_LINK}/${user.user_id}`}
                            alt={user.name}
                        />
                    </div>
                    <div class="user-details">
                        <p class="user-name">{user.name}</p>
                        <p class="user-username">@{user.username}</p>
                    </div>
                </div>
                <button
                    class="unblock-button"
                    onclick={() => handleUnblock(user.user_id, user.username)}
                >
                    Unblock
                </button>
            </div>
        {/each}
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/blocked.scss";
</style>
