<script lang="ts">
    import { onMount } from "svelte";
    import Post from "../components/Post.svelte";
    import type { UserProfile } from "../types/user";
    import { getProfile, getSelfProfile } from "../controllers/user-controller";
    import { addToast } from "../stores/toast-wrapper";
    import { navigate } from "svelte-routing";
    import { BadgeCheck } from "@lucide/svelte";
    import type { Thread } from "../types/thread";

    let user: UserProfile;
    let posts: Thread[] = [];

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

    onMount(async () => {
        await loadProfile();
    });

    const tabs = ["Posts", "Replies", "Likes", "Media"];
    let activeTab = "Posts";

    let showImagePreviewModal = false;
    let previewImageUrl = "";

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
                on:click={() =>
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
                <button
                    class="edit-profile-button"
                    on:click={() => navigate("/settings")}
                >
                    Edit profile
                </button>
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
                        <rect x="3" y="4" width="18" height="18" rx="2" ry="2"
                        ></rect>
                        <line x1="16" y1="2" x2="16" y2="6"></line>
                        <line x1="8" y1="2" x2="8" y2="6"></line>
                        <line x1="3" y1="10" x2="21" y2="10"></line>
                    </svg>
                    Joined {user?.join_date}
                </div>

                <div class="follow-info">
                    <span class="following"
                        ><strong>{user.following}</strong> Following</span
                    >
                    <span class="followers"
                        ><strong>{user.followers}</strong> Followers</span
                    >
                </div>
            </div>
        </div>

        <div class="profile-tabs">
            {#each tabs as tab}
                <button
                    class="tab-button {activeTab === tab ? 'active' : ''}"
                    on:click={() => setActiveTab(tab)}
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

        {#if showImagePreviewModal}
            <div class="modal-overlay" on:click={closeImagePreviewModal}>
                <div class="image-preview-container" on:click|stopPropagation>
                    <button
                        class="close-button"
                        on:click={closeImagePreviewModal}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <line x1="18" y1="6" x2="6" y2="18"></line>
                            <line x1="6" y1="6" x2="18" y2="18"></line>
                        </svg>
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

<style lang="scss">
    @use "../styles/profile.scss";
</style>
