<script lang="ts">
    import { onMount } from "svelte";
    import {
        deleteNotif,
        markNotifAsRead,
    } from "../controllers/notification-controller";
    import { addToast } from "../stores/toast-wrapper";
    import type { Notification } from "../types/notification";
    import type { UserProfile } from "../types/user";
    import { getUserById } from "../controllers/user-controller";

    let { notification }: { notification: Notification } = $props();
    let user = $state<UserProfile | null>(null);

    const formattedTime = notification.timestamp.toString();

    async function handleDelete() {
        const response = await deleteNotif(notification.notification_id);

        if (response.success) {
            addToast("success", "Notification deleted successfully", "Success");
            setTimeout(() => {
                window.location.reload();
            }, 500);
        } else {
            addToast("error", response.message, "Error");
        }
    }

    async function handleMarkAsRead() {
        const response = await markNotifAsRead(notification.notification_id);

        if (response.success) {
            addToast(
                "success",
                "Notification marked as read successfully",
                "Success"
            );
            setTimeout(() => {
                window.location.reload();
            }, 500);
        } else {
            addToast("error", response.message, "Error");
        }
    }

    onMount(async () => {
        if (notification.from === "System") {
            user = {
                name: "System",
                username: "System",
                is_verified: false,
                user_id: notification.from,
                bio: "",
                followers: 0,
                following: 0,
                gender: "",
                date_of_birth: "",
                email: "",
                join_date: "",
                is_private: false,
            };
        }
        const response = await getUserById(notification.from);
        if (response) {
            user = response;
        } else {
            user = {
                name: "Unknown.",
                username: "Unknown.",
                is_verified: false,
                user_id: notification.from,
                bio: "",
                followers: 0,
                following: 0,
                gender: "",
                date_of_birth: "",
                email: "",
                join_date: "",
                is_private: false,
            };
        }
    });

    function handleRedirect() {}
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
    class={`notification-card ${notification.read ? "read" : "unread"}`}
    onclick={() => handleRedirect()}
    onkeydown={(e) => {
        if (e.key === "Enter") {
            handleRedirect();
        }
    }}
>
    <div class="notification-content">
        <h3>{notification.title}</h3>
        <p>{@html notification.content}</p>
        <span class="meta">
            <strong>{user?.username}</strong> &bull; {formattedTime}
        </span>
    </div>

    <div class="actions">
        {#if !notification.read}
            <button class="mark-read" onclick={handleMarkAsRead}
                >Mark as Read</button
            >
        {/if}
        <button class="delete" onclick={handleDelete}>Delete</button>
    </div>
</div>

<style lang="scss">
    @use "../styles/notification-card.scss";
</style>
