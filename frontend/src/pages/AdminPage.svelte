<script lang="ts">
    import { onMount } from "svelte";
    import { isLoggedIn } from "../controllers/token-controller";
    import type { UserProfile, VerifyAccountRequest } from "../types/user";
    import {
        addThreadCategory,
        deleteThreadCategory,
        getAllUserVerificationRequests,
    } from "../controllers/admin-controller";
    import UserVerificationRequest from "../components/Admin/UserVerificationRequest.svelte";
    import { getCategories } from "../controllers/thread-controller";
    import { addToast } from "../stores/toast-wrapper";
    import { Plus, X } from "@lucide/svelte";

    async function loadUserVerificationRequests() {
        const response = await getAllUserVerificationRequests();

        if (response.success && response.payload) {
            verificationRequests = response.payload.requests;
        } else {
            verificationRequests = [];
        }
    }

    async function loadThreadCategories() {
        const response = await getCategories();
        if (response.success && response.payload) {
            threadCategories = response.payload;
        } else {
            threadCategories = [];
            addToast("error", response.message || "Failed to load categories.");
        }
    }

    async function handleAddCategory() {
        if (newThreadCategory.trim() === "") {
            addToast("error", "Category name cannot be empty.");
            return;
        }

        const response = await addThreadCategory(newThreadCategory.trim());
        if (response.success) {
            addToast("success", "Category added successfully.");
            newThreadCategory = "";
            loadThreadCategories();
        } else {
            addToast("error", response.message || "Failed to add category.");
        }
    }

    onMount(async () => {
        if (await !isLoggedIn()) {
            window.location.href = "/login";
        } else if (localStorage.getItem("is_admin") !== "true") {
            window.location.href = "/home";
        }

        loadUserVerificationRequests();
        loadThreadCategories();
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
    let newThreadCategory = $state<string>("");
    let newCommunityCategory = $state<string>("");
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
            <h2>{activeTab}</h2>
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
            {:else if activeTab === "Categories"}
                <div class="categories-container">
                    <div class="categories">
                        <h3>Thread Categories</h3>
                        <div class="category-list">
                            {#each threadCategories as category, i}
                                <div class="category-card">
                                    <span>{category}</span>
                                    <button
                                        class="delete-button"
                                        title="Delete category"
                                        onclick={async () => {
                                            await deleteThreadCategory(
                                                category
                                            );
                                            loadThreadCategories();
                                        }}
                                    >
                                        <X />
                                    </button>
                                </div>
                            {/each}
                        </div>
                        <div class="add-category">
                            <input
                                type="text"
                                placeholder="New Category"
                                bind:value={newThreadCategory}
                            />
                            <button onclick={handleAddCategory}>Add</button>
                        </div>
                    </div>
                    <div class="categories">
                        <h3>Community Categories</h3>
                        <div class="category-list">
                            {#each communityCategories as category, i}
                                <div class="category-card">
                                    <span>{category}</span>
                                    <button
                                        class="delete-button"
                                        title="Delete category"
                                        onclick={async () => {
                                            await deleteThreadCategory(
                                                category
                                            );
                                            loadThreadCategories();
                                        }}
                                    >
                                        <X />
                                    </button>
                                </div>
                            {/each}
                        </div>
                        <div class="add-category">
                            <input
                                type="text"
                                placeholder="New Category"
                                bind:value={newCommunityCategory}
                            />
                            <button>Add</button>
                        </div>
                    </div>
                </div>
            {/if}
        </div>
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/admin.scss";
</style>
