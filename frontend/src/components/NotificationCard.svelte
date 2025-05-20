<script lang="ts">
    import {
        deleteNotif,
        markNotifAsRead,
    } from "../controllers/notification-controller";
    import { addToast } from "../stores/toast-wrapper";
    import type { Notification } from "../types/notification";

    let { notification }: { notification: Notification } = $props();

    const formattedTime = new Date(notification.timestamp).toLocaleString();

    async function handleDelete() {
        const response = await deleteNotif(notification.notification_id);

        if (response.success) {
            addToast("success", "Notification deleted successfully", "Success");
            window.location.reload();
        } else {
            addToast("error", response.message, "Error");
        }
    }

    async function handleMarkAsRead() {
        const response = await markNotifAsRead(notification.notification_id);

        if (response.success) {
            addToast("success", "Notification deleted successfully", "Success");
            window.location.reload();
        } else {
            addToast("error", response.message, "Error");
        }
    }

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
        <p>{notification.content}</p>
        <span class="meta">
            <strong>{notification.from}</strong> &bull; {formattedTime}
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
