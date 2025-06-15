<script lang="ts">
    import { onMount } from "svelte";
    import { isLoggedIn } from "../controllers/token-controller";
    import {
        addThreadCategory,
        deleteThreadCategory,
        getAllUsers,
        getAllUserVerificationRequests,
        sendNewsletter,
        getCommunityRequests,
        approveCommunityRequest,
        rejectCommunityRequest,
        getUserReports,
        approveReport,
        rejectReport,
        addCommunityCategory,
        deleteCommunityCategory,
    } from "../controllers/admin-controller";
    import UserVerificationRequest from "../components/Admin/UserVerificationRequest.svelte";
    import { getCategories } from "../controllers/thread-controller";
    import { addToast } from "../stores/toast-wrapper";
    import { X, Send, Check, AlertTriangle } from "@lucide/svelte";
    import type {
        AdminUserProfile,
        CommunityRequest,
        UserReport,
        VerifyAccountRequest,
    } from "../types/admin";
    import { getCommunityCategories } from "../controllers/community-controller";
    import AdminUserView from "../components/Admin/AdminUserView.svelte";
    import AdminReportView from "../components/Admin/AdminReportView.svelte";

    let activeTab = $state<string>("Users");
    let users = $state<AdminUserProfile[]>([]);
    let reports = $state<UserReport[]>([]);
    let newsLetterTitle = $state<string>("");
    let newsLetterContent = $state<string>("");
    let communityRequests = $state<CommunityRequest[]>([]);
    let verificationRequests = $state<VerifyAccountRequest[]>([]);
    let threadCategories = $state<string[]>([]);
    let communityCategories = $state<string[]>([]);
    let newThreadCategory = $state<string>("");
    let newCommunityCategory = $state<string>("");
    let isLoading = $state<boolean>(false);

    const tabs = [
        "Users",
        "Newsletter",
        "Community",
        "Premium",
        "Reports",
        "Categories",
    ];

    async function loadUserVerificationRequests() {
        try {
            const response = await getAllUserVerificationRequests();
            if (response.success && response.payload) {
                verificationRequests = response.payload.requests;
            } else {
                verificationRequests = [];
            }
        } catch (error) {
            addToast("error", "Failed to load verification requests");
        }
    }

    async function loadThreadCategories() {
        try {
            const response = await getCategories();
            if (response.success && response.payload) {
                threadCategories = response.payload;
            } else {
                threadCategories = [];
                addToast("error", "Failed to load thread categories");
            }
        } catch (error) {
            addToast("error", "Failed to load thread categories");
        }
    }

    async function loadCommunityCategories() {
        try {
            const response = await getCommunityCategories();
            if (response.success && response.payload) {
                communityCategories = response.payload;
            } else {
                communityCategories = [];
            }
        } catch (error) {
            addToast("error", "Failed to load community categories");
        }
    }

    async function loadAllUsers() {
        try {
            const response = await getAllUsers();
            if (response.success && response.payload) {
                users = response.payload.users;
            } else {
                users = [];
                addToast("error", "Failed to load users");
            }
        } catch (error) {
            addToast("error", "Failed to load users");
        }
    }

    async function loadCommunityRequests() {
        try {
            const response = await getCommunityRequests();
            if (response.success && response.payload) {
                communityRequests = response.payload.requests;
            } else {
                communityRequests = [];
            }
        } catch (error) {
            addToast("error", "Failed to load community requests");
        }
    }

    async function loadUserReports() {
        try {
            const response = await getUserReports();
            if (response.success && response.payload) {
                reports = response.payload.reports;
            } else {
                reports = [];
            }
        } catch (error) {
            addToast("error", "Failed to load user reports");
        }
    }

    async function handleSendNewsletter() {
        if (!newsLetterTitle.trim() || !newsLetterContent.trim()) {
            addToast("error", "Please fill in both title and content");
            return;
        }

        try {
            isLoading = true;
            const response = await sendNewsletter(
                newsLetterTitle,
                newsLetterContent
            );
            if (response.success) {
                addToast("success", "Newsletter sent successfully");
                newsLetterTitle = "";
                newsLetterContent = "";
            } else {
                addToast(
                    "error",
                    response.message || "Failed to send newsletter"
                );
            }
        } catch (error) {
            addToast("error", "Failed to send newsletter");
        } finally {
            isLoading = false;
        }
    }

    async function handleApproveCommunityRequest(requestId: string) {
        try {
            const response = await approveCommunityRequest(requestId);
            if (response.success) {
                addToast("success", "Community request approved");
                await loadCommunityRequests();
            } else {
                addToast(
                    "error",
                    response.message || "Failed to approve request"
                );
            }
        } catch (error) {
            addToast("error", "Failed to approve request");
        }
    }

    async function handleRejectCommunityRequest(requestId: string) {
        try {
            const response = await rejectCommunityRequest(requestId);
            if (response.success) {
                addToast("success", "Community request rejected");
                await loadCommunityRequests();
            } else {
                addToast(
                    "error",
                    response.message || "Failed to reject request"
                );
            }
        } catch (error) {
            addToast("error", "Failed to reject request");
        }
    }

    async function handleAddThreadCategory() {
        if (newThreadCategory.trim() === "") {
            addToast("error", "Category name cannot be empty");
            return;
        }

        try {
            const response = await addThreadCategory(newThreadCategory.trim());
            if (response.success) {
                addToast("success", "Thread category added successfully");
                newThreadCategory = "";
                await loadThreadCategories();
            } else {
                addToast("error", response.message || "Failed to add category");
            }
        } catch (error) {
            addToast("error", "Failed to add category");
        }
    }

    async function handleAddCommunityCategory() {
        if (newCommunityCategory.trim() === "") {
            addToast("error", "Category name cannot be empty");
            return;
        }

        try {
            const response = await addCommunityCategory(
                newCommunityCategory.trim()
            );
            if (response.success) {
                addToast("success", "Community category added successfully");
                newCommunityCategory = "";
                await loadCommunityCategories();
            } else {
                addToast("error", response.message || "Failed to add category");
            }
        } catch (error) {
            addToast("error", "Failed to add category");
        }
    }

    async function handleDeleteThreadCategory(category: string) {
        try {
            const response = await deleteThreadCategory(category);
            if (response.success) {
                addToast("success", "Thread category deleted successfully");
                await loadThreadCategories();
            } else {
                addToast(
                    "error",
                    response.message || "Failed to delete category"
                );
            }
        } catch (error) {
            addToast("error", "Failed to delete category");
        }
    }

    async function handleDeleteCommunityCategory(category: string) {
        try {
            const response = await deleteCommunityCategory(category);
            if (response.success) {
                addToast("success", "Community category deleted successfully");
                await loadCommunityCategories();
            } else {
                addToast(
                    "error",
                    response.message || "Failed to delete category"
                );
            }
        } catch (error) {
            addToast("error", "Failed to delete category");
        }
    }

    onMount(async () => {
        if (!(await isLoggedIn())) {
            window.location.href = "/login";
            return;
        }

        if (localStorage.getItem("is_admin") !== "true") {
            window.location.href = "/home";
            return;
        }

        await Promise.all([
            loadAllUsers(),
            loadUserVerificationRequests(),
            loadThreadCategories(),
            loadCommunityCategories(),
            loadCommunityRequests(),
            loadUserReports(),
        ]);
    });
</script>

<div class="admin-container">
    <div class="admin-header">
        <h1>Admin Dashboard</h1>
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
                <div class="users-container">
                    <h3>User Management ({users.length} users)</h3>
                    <div class="user-list">
                        {#each users as user}
                            <AdminUserView {user} />
                        {/each}
                    </div>
                </div>
            {:else if activeTab === "Newsletter"}
                <div class="newsletter-container">
                    <h3>Send Newsletter</h3>
                    <div class="newsletter-form">
                        <div class="form-group">
                            <label for="newsletter-title">Title</label>
                            <input
                                id="newsletter-title"
                                type="text"
                                bind:value={newsLetterTitle}
                                placeholder="Newsletter title..."
                            />
                        </div>
                        <div class="form-group">
                            <label for="newsletter-content">Content</label>
                            <textarea
                                id="newsletter-content"
                                bind:value={newsLetterContent}
                                placeholder="Newsletter content..."
                                rows="10"
                            ></textarea>
                        </div>
                        <button
                            class="action-btn send-btn"
                            onclick={handleSendNewsletter}
                            disabled={isLoading ||
                                !newsLetterTitle.trim() ||
                                !newsLetterContent.trim()}
                        >
                            <Send size={16} />
                            Send Newsletter
                        </button>
                    </div>
                </div>
            {:else if activeTab === "Community"}
                <div class="community-container">
                    <h3>
                        Community Requests ({communityRequests.length} pending)
                    </h3>
                    <div class="request-list">
                        {#each communityRequests as request}
                            <div class="request-card">
                                <div class="request-header">
                                    <h4>{request.community_name}</h4>
                                    <span class="timestamp"
                                        >{new Date(
                                            request.submitted_at
                                        ).toLocaleDateString()}</span
                                    >
                                </div>
                                <div class="request-info">
                                    <p>
                                        <strong>Category:</strong>
                                        {request.category}
                                    </p>
                                    <p>
                                        <strong>Description:</strong>
                                        {request.description}
                                    </p>
                                    <p>
                                        <strong>Requested by:</strong>
                                        {request.user_id}
                                    </p>
                                </div>
                                <!-- {#if request.community_image_url}
                                    <img
                                        src={request.community_image_url ||
                                            "/placeholder.svg"}
                                        alt="Community"
                                        class="community-preview"
                                    />
                                {/if} -->
                                <div class="actions">
                                    <button
                                        class="action-btn approve-btn"
                                        onclick={() =>
                                            handleApproveCommunityRequest(
                                                request.community_id
                                            )}
                                    >
                                        <Check size={16} />
                                        Approve
                                    </button>
                                    <button
                                        class="action-btn reject-btn"
                                        onclick={() =>
                                            handleRejectCommunityRequest(
                                                request.community_id
                                            )}
                                    >
                                        <X size={16} />
                                        Reject
                                    </button>
                                </div>
                            </div>
                        {/each}
                        {#if communityRequests.length === 0}
                            <p class="empty-state">
                                No pending community requests
                            </p>
                        {/if}
                    </div>
                </div>
            {:else if activeTab === "Premium"}
                <div class="premium-container">
                    <h3>
                        Premium Verification Requests ({verificationRequests.length}
                        pending)
                    </h3>
                    {#each verificationRequests as request}
                        <UserVerificationRequest {request} />
                    {/each}
                    {#if verificationRequests.length === 0}
                        <p class="empty-state">
                            No pending verification requests
                        </p>
                    {/if}
                </div>
            {:else if activeTab === "Reports"}
                <div class="reports-container">
                    <h3>
                        User Reports ({reports.filter(
                            (r) => r.status === "pending"
                        ).length} pending)
                    </h3>
                    <div class="report-list">
                        {#each reports as report}
                            <AdminReportView {report} />
                        {/each}
                        {#if reports.length === 0}
                            <p class="empty-state">No pending reports</p>
                        {/if}
                    </div>
                </div>
            {:else if activeTab === "Categories"}
                <div class="categories-container">
                    <div class="categories">
                        <h3>Thread Categories</h3>
                        <div class="category-list">
                            {#each threadCategories as category}
                                <div class="category-card">
                                    <span class="category-name">{category}</span
                                    >
                                    <button
                                        class="delete-button"
                                        title="Delete category"
                                        onclick={() =>
                                            handleDeleteThreadCategory(
                                                category
                                            )}
                                    >
                                        <X size={16} />
                                    </button>
                                </div>
                            {/each}
                        </div>
                        <div class="add-category">
                            <input
                                type="text"
                                placeholder="New thread category"
                                bind:value={newThreadCategory}
                            />
                            <button onclick={handleAddThreadCategory}
                                >Add</button
                            >
                        </div>
                    </div>

                    <div class="categories">
                        <h3>Community Categories</h3>
                        <div class="category-list">
                            {#each communityCategories as category}
                                <div class="category-card">
                                    <span class="category-name">{category}</span
                                    >
                                    <button
                                        class="delete-button"
                                        title="Delete category"
                                        onclick={() =>
                                            handleDeleteCommunityCategory(
                                                category
                                            )}
                                    >
                                        <X size={16} />
                                    </button>
                                </div>
                            {/each}
                        </div>
                        <div class="add-category">
                            <input
                                type="text"
                                placeholder="New community category"
                                bind:value={newCommunityCategory}
                            />
                            <button onclick={handleAddCommunityCategory}
                                >Add</button
                            >
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
