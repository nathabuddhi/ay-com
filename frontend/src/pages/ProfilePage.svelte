<script lang="ts">
    import { onMount } from "svelte";
    import Post from "../components/Post.svelte";
    import type { UserProfile } from "../types/user";
    import {
        blockUser,
        followUser,
        getProfile,
        getSelfProfile,
        unfollowUser,
    } from "../controllers/user-controller";
    import { addToast } from "../stores/toast-wrapper";
    import { navigate } from "svelte-routing";
    import { BadgeCheck, Calendar, CrossIcon, X } from "@lucide/svelte";
    import type { Thread } from "../types/thread";

    let user = $state<UserProfile>({
        name: "loading",
        username: "loading",
        is_verified: false,
        user_id: "loading",
        bio: "",
        followers: 0,
        following: 0,
        gender: "",
        date_of_birth: "",
        email: "",
        join_date: "",
    });
    let posts = $state<Thread[]>([]);
    let replies = $state<Thread[]>([]);
    let likes = $state<Thread[]>([]);
    let media = $state<Thread[]>([]);

    const tabs = ["Posts", "Replies", "Likes", "Media"];
    let activeTab = $state("Posts");

    let showImagePreviewModal = $state(false);
    let previewImageUrl = $state("");

    async function loadProfile() {
        const username = window.location.pathname.split("/").pop();
        let response;
        if (username === "profile") {
            response = await getSelfProfile().catch((error) => {
                addToast(
                    "error",
                    "Error fetching user profile" + error,
                    "Error!"
                );
            });
        } else {
            response = await getProfile(username ?? "").catch((error) => {
                addToast(
                    "error",
                    "Error fetching user profile" + error,
                    "Error!"
                );
            });
        }

        if (response && response.success && response.payload) {
            user = response.payload;
        } else {
            addToast(
                "error",
                "Error fetching user profile: " + response?.message,
                "Error!"
            );
        }
    }

    async function handleBlock(userId: string) {
        const response = await blockUser(userId);

        if (response.success) {
            addToast("success", "Blocked successfully!", "Success!");
        } else {
            addToast("error", response.message, "Error!");
        }
    }

    async function handleUnfollow(userId: string) {
        const response = await unfollowUser(userId);

        if (response.success) {
            addToast("success", "Unfollowed successfully!", "Success!");
        } else {
            addToast("error", response.message, "Error!");
        }
    }

    async function handleFollow(userId: string) {
        const response = await followUser(userId);

        if (response.success) {
            addToast("success", "Followed successfully!", "Success!");
        } else {
            addToast("error", response.message, "Error!");
        }
    }

    onMount(async () => {
        await loadProfile();
    });

    function setActiveTab(tab: string) {
        activeTab = tab;
    }

    function openImagePreviewModal(imageUrl: string) {
        previewImageUrl = imageUrl;
        showImagePreviewModal = true;
    }

    function closeImagePreviewModal() {
        showImagePreviewModal = false;
    }

    function navigateToThread(id: string) {
        console.log(`Navigating to thread ${id}`);
    }
</script>

{#if user}
    <div class="profile-container">
        <div class="profile-header">
            <div class="header-info">
                <h1>{user?.name || "Loading"}</h1>
                <span class="post-count">
                    {posts.length} {posts.length === 1 ? "Post" : "Posts"}</span
                >
            </div>
        </div>

        <div class="profile-banner">
            <!-- <img
                src={import.meta.env.VITE_BANNER_LINK + user?.user_id + ".png"}
                alt="Banner"
            /> -->
            <div
                class="profile-avatar"
                onclick={() =>
                    openImagePreviewModal(
                        import.meta.env.VITE_AVATAR_LINK +
                            user?.user_id +
                            ".png"
                    )}
            >
                <img
                    src={import.meta.env.VITE_AVATAR_LINK +
                        user?.user_id +
                        ".png"}
                    alt={user?.name || "Loading"}
                />
            </div>
        </div>
        <div class="profile-info">
            <div class="profile-actions">
                {#if user.username === localStorage.getItem("username")}
                    <button
                        class="edit-profile-button"
                        onclick={() => navigate("/settings")}
                    >
                        Edit profile
                    </button>
                {:else}
                    <button
                        class="edit-profile-button"
                        onclick={() => {
                            handleBlock(user.user_id);
                        }}
                    >
                        Block
                    </button>
                    <button
                        class="edit-profile-button"
                        onclick={() => {
                            handleFollow(user.user_id);
                        }}
                    >
                        Follow
                    </button>
                    <button
                        class="edit-profile-button"
                        onclick={() => {
                            handleUnfollow(user.user_id);
                        }}
                    >
                        Unfollow
                    </button>
                {/if}
            </div>

            <div class="profile-name-section">
                <h2 class="profile-name">{user?.name || "Loading"}</h2>
                {#if user?.is_verified || false}
                    <span class="verified-badge">
                        <BadgeCheck />
                        Get verified
                    </span>
                {/if}
            </div>

            <div class="profile-username">@{user?.username || "Loading"}</div>

            {#if user?.bio || false}
                <div class="profile-bio">{user.bio}</div>
            {/if}

            <div class="profile-meta">
                <div class="joined-date">
                    <Calendar />
                    Joined {user?.join_date}
                </div>

                <div class="follow-info">
                    <span class="following"
                        ><strong>{user.following ? user.following : 0}</strong> Following</span
                    >
                    <span class="followers"
                        ><strong>{user.followers ? user.followers : 0}</strong> Followers</span
                    >
                </div>
            </div>
        </div>

        <div class="profile-tabs">
            {#each tabs as tab}
                <button
                    class="tab-button {activeTab === tab ? 'active' : ''}"
                    onclick={() => setActiveTab(tab)}
                >
                    {tab}
                </button>
            {/each}
        </div>

        <!-- <div class="profile-content">
            {#if activeTab === "Posts"}
                <div class="posts-container">
                    {#if posts.length === 0}
                        <div class="empty-state">No posts yet</div>
                    {:else}
                        {#each posts.filter((post) => post.pinned) as post (post.thread_id)}
                            <div class="post-item">
                                <div class="post-header">
                                    <div class="post-user-info">
                                        <span class="post-pinned">
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                width="16"
                                                height="16"
                                                viewBox="0 0 24 24"
                                                fill="none"
                                                stroke="currentColor"
                                                stroke-width="2"
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                            >
                                                <path d="M12 2L12 12"></path>
                                                <path d="M12 22L12 12"></path>
                                                <path d="M4.93 10.93L19.07 10.93"
                                                ></path>
                                                <path d="M4.93 13.07L19.07 13.07"
                                                ></path>
                                            </svg>
                                            Pinned
                                        </span>
                                    </div>
                                    <div class="post-actions-dropdown">
                                        <button class="post-more-options">
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                width="16"
                                                height="16"
                                                viewBox="0 0 24 24"
                                                fill="none"
                                                stroke="currentColor"
                                                stroke-width="2"
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                            >
                                                <circle cx="12" cy="12" r="1"
                                                ></circle>
                                                <circle cx="12" cy="5" r="1"
                                                ></circle>
                                                <circle cx="12" cy="19" r="1"
                                                ></circle>
                                            </svg>
                                        </button>
                                    </div>
                                </div>
                                <div
                                    class="post-content"
                                    on:click={() => navigateToThread(post.thread_id)}
                                >
                                    <Post {post} />
                                </div>
                            </div>
                        {/each}

                        {#each posts.filter((post) => !post.pinned) as post (post.thread_id)}
                            <div class="post-item">
                                <div
                                    class="post-content"
                                    on:click={() => navigateToThread(post.thread_id)}
                                >
                                    <Post {post} />
                                </div>
                                <div class="post-actions-dropdown">
                                    <button
                                        class="post-more-options"
                                        on:click={() => pinItem("post", post.thread_id)}
                                    >
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            width="16"
                                            height="16"
                                            viewBox="0 0 24 24"
                                            fill="none"
                                            stroke="currentColor"
                                            stroke-width="2"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        >
                                            <circle cx="12" cy="12" r="1"></circle>
                                            <circle cx="12" cy="5" r="1"></circle>
                                            <circle cx="12" cy="19" r="1"></circle>
                                        </svg>
                                    </button>
                                </div>
                            </div>
                        {/each}
                    {/if}
                </div>
            {:else if activeTab === "Replies"}
                <div class="replies-container">
                    {#if replies.length === 0}
                        <div class="empty-state">No replies yet</div>
                    {:else}
                        {#each replies.filter((reply) => reply.isPinned) as reply (reply.id)}
                            <div class="reply-item">
                                <div class="post-header">
                                    <div class="post-user-info">
                                        <span class="post-pinned">
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                width="16"
                                                height="16"
                                                viewBox="0 0 24 24"
                                                fill="none"
                                                stroke="currentColor"
                                                stroke-width="2"
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                            >
                                                <path d="M12 2L12 12"></path>
                                                <path d="M12 22L12 12"></path>
                                                <path d="M4.93 10.93L19.07 10.93"
                                                ></path>
                                                <path d="M4.93 13.07L19.07 13.07"
                                                ></path>
                                            </svg>
                                            Pinned
                                        </span>
                                    </div>
                                    <div class="post-actions-dropdown">
                                        <button class="post-more-options">
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                width="16"
                                                height="16"
                                                viewBox="0 0 24 24"
                                                fill="none"
                                                stroke="currentColor"
                                                stroke-width="2"
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                            >
                                                <circle cx="12" cy="12" r="1"
                                                ></circle>
                                                <circle cx="12" cy="5" r="1"
                                                ></circle>
                                                <circle cx="12" cy="19" r="1"
                                                ></circle>
                                            </svg>
                                        </button>
                                    </div>
                                </div>
                                <div
                                    class="reply-content"
                                    on:click={() =>
                                        navigateToThread(reply.parentPostId)}
                                >
                                    <Post post={reply} />
                                </div>
                            </div>
                        {/each}

                        {#each replies.filter((reply) => !reply.isPinned) as reply (reply.id)}
                            <div class="reply-item">
                                <div
                                    class="reply-content"
                                    on:click={() =>
                                        navigateToThread(reply.parentPostId)}
                                >
                                    <Post post={reply} />
                                </div>
                                <div class="post-actions-dropdown">
                                    <button
                                        class="post-more-options"
                                        on:click={() => pinItem("reply", reply.id)}
                                    >
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            width="16"
                                            height="16"
                                            viewBox="0 0 24 24"
                                            fill="none"
                                            stroke="currentColor"
                                            stroke-width="2"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        >
                                            <circle cx="12" cy="12" r="1"></circle>
                                            <circle cx="12" cy="5" r="1"></circle>
                                            <circle cx="12" cy="19" r="1"></circle>
                                        </svg>
                                    </button>
                                </div>
                            </div>
                        {/each}
                    {/if}
                </div>
            {:else if activeTab === "Likes"}
                <div class="likes-container">
                    {#if likedPosts.length === 0}
                        <div class="empty-state">No likes yet</div>
                    {:else}
                        {#each likedPosts as post (post.thread_id)}
                            <div class="post-item">
                                <div
                                    class="post-content"
                                    on:click={() => navigateToThread(post.thread_id)}
                                >
                                    <Post {post} />
                                </div>
                            </div>
                        {/each}
                    {/if}
                </div>
            {:else if activeTab === "Media"}
                <div class="media-container">
                    {#if mediaItems.length === 0}
                        <div class="empty-state">No media yet</div>
                    {:else}
                        <div class="media-grid">
                            {#each mediaItems as media (media.id)}
                                <div
                                    class="media-item"
                                    on:click={() => navigateToThread(media.postId)}
                                >
                                    {#if media.type === "image" || media.type === "gif"}
                                        <img
                                            src={media.url || "/placeholder.svg"}
                                            alt="Media content"
                                        />
                                    {:else if media.type === "video"}
                                        <video src={media.url} controls={false}
                                        ></video>
                                        <div class="video-overlay">
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                width="24"
                                                height="24"
                                                viewBox="0 0 24 24"
                                                fill="white"
                                                stroke="white"
                                                stroke-width="2"
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                            >
                                                <polygon points="5 3 19 12 5 21 5 3"
                                                ></polygon>
                                            </svg>
                                        </div>
                                    {/if}
                                </div>
                            {/each}
                        </div>
                    {/if}
                </div>
            {/if}
        </div> -->

        <!-- svelte-ignore a11y_click_events_have_key_events -->
        {#if showImagePreviewModal}
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div class="modal-overlay" onclick={closeImagePreviewModal}>
                <div class="image-preview-container">
                    <!-- svelte-ignore a11y_consider_explicit_label -->
                    <button
                        class="close-button"
                        onclick={closeImagePreviewModal}
                    >
                        <X />
                    </button>
                    <img
                        src={previewImageUrl || "/placeholder.svg"}
                        alt="Preview"
                    />
                </div>
            </div>
        {/if}
    </div>
{/if}

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/profile.scss";
</style>
