<script lang="ts">
    import { onMount } from "svelte";
    import { isLoggedIn } from "../controllers/token-controller";

    onMount(async () => {
        if (await !isLoggedIn()) {
            window.location.href = "/login";
        } else if (localStorage.getItem("is_admin") !== "true") {
            window.location.href = "/home";
        }
    });

    const tabs = [
        "Users",
        "NewsLetter",
        "Community",
        "Premium",
        "Reports",
        "Categories",
    ];

    let activeTab = $state<string>("Users");
</script>

<div class="admin-container">
    <div class="admin-header">
        <h1>Admin</h1>
        <div class="admin-tabs">
            {#each tabs as tab}
                <button
                    class="tab-button {activeTab === tab ? 'active' : ''}"
                    onclick={() => (activeTab = tab)}
                >
                    {tab}
                </button>
            {/each}
        </div>
    </div>
</div>

<style lang="scss">
    @use "../styles/admin.scss";
</style>
