<script lang="ts">
    import { onMount } from "svelte";
    import { isLoggedIn } from "../controllers/token-controller";
    import type { UserProfile, VerifyAccountRequest } from "../types/user";
    import { getAllUserVerificationRequests } from "../controllers/admin-controller";
    import UserVerificationRequest from "../components/Admin/UserVerificationRequest.svelte";

    async function loadUserVerificationRequests() {
        const response = await getAllUserVerificationRequests();

        if (response.success && response.payload) {
            verificationRequests = response.payload.requests;
        } else {
            verificationRequests = [];
        }
    }

    onMount(async () => {
        if (await !isLoggedIn()) {
            window.location.href = "/login";
        } else if (localStorage.getItem("is_admin") !== "true") {
            window.location.href = "/home";
        }

        loadUserVerificationRequests();
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

    let users = $state<UserProfile[]>([]);

    let newsLetterTitle = $state<string>("");
    let newsLetterContent = $state<string>("");

    // let communityRequests = $state<CommunityRequest[]>([]);
    let verificationRequests = $state<VerifyAccountRequest[]>([]);

    let threadCategories = $state<string[]>([]);
    let communityCategories = $state<string[]>([]);
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
    <div class="admin-content">
        <div class="tab-content">
            <h3>{activeTab}</h3>
            {#if activeTab === "Users"}
                <ul>
                    {#each users as user}
                        <li>{user.username} - {user.email}</li>
                    {/each}
                </ul>
            {:else if activeTab === "Premium"}
                {#each verificationRequests as request}
                    <UserVerificationRequest {request} />
                {/each}
                {#if verificationRequests.length === 0}
                    <p>No pending verification requests found.</p>
                {/if}
            {/if}
        </div>
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/admin.scss";
</style>
