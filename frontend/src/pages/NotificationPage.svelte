<script lang="ts">
    import { onMount } from "svelte";
    import { addToast } from "../stores/toast-wrapper";
    import type { Notification } from "../types/notification";
    import {
        clearNotifs,
        getAllNotifications,
    } from "../controllers/notification-controller";
    import NotificationCard from "../components/NotificationCard.svelte";

    let notifications = $state<Notification[]>([]);
    let mentions = $state<Notification[]>([]);
    const tabs = ["All", "Mentions"];
    let activeTab = $state("All");

    async function fetchNotifications() {
        const response = await getAllNotifications();
        if (response.success && response.payload) {
            notifications = response.payload.notifications;
            mentions = notifications.filter((notif) =>
                notif.title.includes("Mentioned")
            );
        } else if (!response.success) {
            addToast("error", response.message, "Error");
        }
    }

    async function handleClearNotif() {
        const response = await clearNotifs();
        if (response.success) {
            addToast(
                "success",
                "Notifications cleared successfully",
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
        await fetchNotifications();
    });

    function setActiveTab(tab: string) {
        activeTab = tab;
    }
</script>

<div class="notification-container">
    <div class="notification-header">
        <div class="header-info">
            <h1>Your Notifications</h1>
            <span class="post-count">
                You have {notifications.filter((n) => !n.read).length ?? 0}
                unread notification{notifications.filter((n) => !n.read)
                    .length > 1
                    ? "s"
                    : ""}
            </span>
        </div>
    </div>
    <div class="notification-actions">
        <button class="delete-all-button" onclick={() => handleClearNotif()}>
            Delete All Notifications
        </button>
    </div>
    <div class="notification-tabs">
        {#each tabs as tab}
            <button
                class="tab-button {activeTab === tab ? 'active' : ''}"
                onclick={() => setActiveTab(tab)}
            >
                {tab}
            </button>
        {/each}
    </div>
    <div class="notification-list">
        {#if activeTab === "All"}
            {#each notifications as notif}
                <NotificationCard notification={notif} />
            {/each}
        {:else}
            {#each mentions as notif}
                <NotificationCard notification={notif} />
            {/each}
        {/if}
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/notifications.scss";
</style>
