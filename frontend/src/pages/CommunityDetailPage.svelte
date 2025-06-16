<script lang="ts">
    import { onMount } from "svelte";
    import type { Community, CommunityMember } from "../types/community";
    import {
        getCommunityById,
        getCommunityMembers,
        // getCommunityMedia,
        sendJoinRequest,
        acceptJoinRequest,
        declineJoinRequest,
        promoteMember,
        demoteModerator,
        getCommunityJoinRequests,
    } from "../controllers/community-controller";
    import { getCommunityThreads } from "../controllers/thread-controller";
    import Post from "../components/Post.svelte";
    import MembersModal from "../components/MembersModal.svelte";
    import {
        AVATAR_IMG,
        C_BANNER_IMG,
        C_ICON_IMG,
        THREAD_IMG,
    } from "../env_var";
    import type { UserProfile } from "../types/user";
    import type { Thread } from "../types/thread";
    import { getUserById } from "../controllers/user-controller";
    import { addToast } from "../stores/toast-wrapper";

    let communityId = $state("");

    let community = $state<Community | null>(null);
    let activeTab = $state<"top" | "latest" | "media" | "about" | "manage">(
        "top"
    );
    let loading = $state(true);
    let showMembersModal = $state(false);

    let topMembers = $state<UserProfile[]>([]);
    let topThreads = $state<Thread[]>([]);
    let latestThreads = $state<Thread[]>([]);
    let allMembers = $state<CommunityMember[]>([]);
    let members = $state<CommunityMember[]>([]);
    let moderators = $state<CommunityMember[]>([]);
    let userProfiles = $state<Record<string, UserProfile>>({});
    let pendingRequests = $state<CommunityMember[]>([]);
    let mediaThreads = $state<Thread[]>([]);

    onMount(async () => {
        communityId = window.location.href.split("/").pop() || "";
        if (communityId) {
            await loadCommunityData();
        }
    });

    async function loadCommunityData() {
        loading = true;
        try {
            const communityRes = await getCommunityById(communityId);
            if (communityRes && communityRes.success && communityRes.payload) {
                community = communityRes.payload;
            }

            const membersRes = await getCommunityMembers(communityId);
            if (membersRes.success && membersRes.payload) {
                allMembers = membersRes.payload.members;
                members = allMembers.filter((m) => m.role === "memb er");
                moderators = allMembers.filter((m) => m.role === "moderator");

                userProfiles = await Promise.all(
                    allMembers.map(async (member) => {
                        const user = await getUserById(member.user_id);
                        return {
                            user_id: member.user_id,
                            user,
                        };
                    })
                ).then((profiles) =>
                    profiles.reduce(
                        (acc, curr) => {
                            if (curr.user) {
                                acc[curr.user_id] = curr.user;
                            }
                            return acc;
                        },
                        {} as Record<string, UserProfile>
                    )
                );
            }

            await loadTabData();
        } finally {
            loading = false;
        }
    }

    async function loadThreads() {
        if (!community) return;

        const response = await getCommunityThreads(community.community_id);
        if (response.success && response.payload) {
            latestThreads = response.payload.threads;
            topThreads = response.payload.threads
                .sort((a, b) => (b.like_count ?? 0) - (a.like_count ?? 0))
                .slice(0, 3);
        } else {
            addToast("error", response.message || "Failed to load threads.");
        }
    }

    async function loadMedia() {
        if (!community) return;

        const response = await getCommunityThreads(community.community_id);
        if (response.success && response.payload) {
            mediaThreads = response.payload.threads.filter(
                (thread) => thread.media && thread.media.length > 0
            );
        } else {
            addToast(
                "error",
                response.message || "Failed to load media threads."
            );
        }
    }

    async function loadTabData() {
        if (!community) return;

        try {
            switch (activeTab) {
                case "top":
                    topMembers = allMembers
                        .map((m) => userProfiles[m.user_id])
                        .filter(Boolean)
                        .sort((a, b) => (b.followers ?? 0) - (a.followers ?? 0))
                        .slice(0, 3);

                    loadThreads();

                    break;

                case "latest":
                    loadThreads();
                    break;

                case "media":
                    loadMedia();
                    break;

                case "manage":
                    if (
                        community.role === "moderator" ||
                        community.role === "owner"
                    ) {
                        const response = await getCommunityMembers(
                            community.community_id
                        );

                        if (response.success && response.payload) {
                            members = response.payload.members.filter(
                                (m) => m.role === "member"
                            );
                            moderators = response.payload.members.filter(
                                (m) => m.role === "moderator"
                            );

                            userProfiles = await Promise.all(
                                response.payload.members.map(async (member) => {
                                    const user = await getUserById(
                                        member.user_id
                                    );
                                    return {
                                        user_id: member.user_id,
                                        user,
                                    };
                                })
                            ).then((profiles) =>
                                profiles.reduce(
                                    (acc, curr) => {
                                        if (curr.user) {
                                            acc[curr.user_id] = curr.user;
                                        }
                                        return acc;
                                    },
                                    {} as Record<string, UserProfile>
                                )
                            );
                        }

                        const pendingResponse = await getCommunityJoinRequests(
                            community.community_id
                        );

                        if (
                            pendingResponse.success &&
                            pendingResponse.payload
                        ) {
                            pendingRequests = pendingResponse.payload.members;

                            const pendingProfiles = await Promise.all(
                                pendingRequests.map(async (member) => {
                                    const user = await getUserById(
                                        member.user_id
                                    );
                                    return {
                                        user_id: member.user_id,
                                        user,
                                    };
                                })
                            ).then((profiles) =>
                                profiles.reduce(
                                    (acc, curr) => {
                                        if (curr.user) {
                                            acc[curr.user_id] = curr.user;
                                        }
                                        return acc;
                                    },
                                    {} as Record<string, UserProfile>
                                )
                            );
                            userProfiles = {
                                ...userProfiles,
                                ...pendingProfiles,
                            };
                        } else {
                            addToast(
                                "error",
                                pendingResponse.message ||
                                    "Failed to load pending requests.",
                                "Failed to load user join requests."
                            );
                            pendingRequests = [];
                        }
                    }
                    break;
            }
        } catch (error) {
            console.error("Error loading tab data:", error);
        }
    }

    async function handleTabChange(tab: typeof activeTab) {
        activeTab = tab;
        await loadTabData();
    }

    async function handleJoinRequest(e: MouseEvent) {
        e.stopPropagation();

        if (!community) return;
        const response = await sendJoinRequest(community.community_id);
        if (response.success) {
            community.is_pending = true;
            addToast("success", "Sent join request succesfully!");
        } else {
            addToast("error", response.message);
        }
    }

    async function handleAcceptRequest(userId: string) {
        if (!community) return;
        try {
            await acceptJoinRequest(community.community_id, userId);
            await loadTabData();
        } catch (error) {
            console.error("Error accepting request:", error);
        }
    }

    async function handleRejectRequest(userId: string) {
        if (!community) return;
        try {
            await declineJoinRequest(community.community_id, userId);
            await loadTabData();
        } catch (error) {
            console.error("Error rejecting request:", error);
        }
    }

    async function handlePromoteMember(userId: string) {
        if (!community) return;
        try {
            await promoteMember(community.community_id, userId);
            await loadTabData();
        } catch (error) {
            console.error("Error promoting member:", error);
        }
    }

    async function handleDemoteModerator(userId: string) {
        if (!community) return;
        try {
            await demoteModerator(community.community_id, userId);
            await loadTabData();
        } catch (error) {
            console.error("Error demoting moderator:", error);
        }
    }

    function navigateToProfile(username: string) {
        window.location.href = `/profile/${username}`;
    }

    function navigateToThread(threadId: string) {
        window.location.href = `/thread/${threadId}`;
    }

    $effect(() => {
        $inspect("CommunityDetailPage", {
            communityId,
            community,
            activeTab,
            loading,
            showMembersModal,
            topMembers,
            topThreads,
            latestThreads,
            allMembers,
            members,
            moderators,
            userProfiles,
            pendingRequests,
            mediaThreads,
        });
    });
</script>

{#if loading}
    <div class="loading">Loading community...</div>
{:else if community}
    <div class="community-detail">
        <div class="community-header">
            {#if community.banner_image}
                <img
                    src={`${C_BANNER_IMG}${community.banner_image}` ||
                        "/placeholder.svg"}
                    alt="Community Banner"
                    class="banner"
                />
            {/if}

            <div class="community-info">
                <div class="community-main">
                    <img
                        src={`${C_ICON_IMG}${community.icon_image}` ||
                            "/placeholder.svg"}
                        alt={community.community_name}
                        class="logo"
                    />
                    <div class="details">
                        <h1>{community.community_name}</h1>
                        <p class="description">{community.description}</p>
                        <div class="categories">
                            {#each community.categories as cat}
                                <span class="category-tag">{cat}</span>
                            {/each}
                        </div>
                    </div>
                </div>

                <div class="community-actions">
                    <button
                        class="members-btn"
                        onclick={() => (showMembersModal = true)}
                    >
                        {community.member_count} Members
                    </button>

                    {#if community.is_pending}
                        <button class="pending-btn" disabled>
                            Request Pending
                        </button>
                    {:else if community.role === "guest"}
                        <button class="join-btn" onclick={handleJoinRequest}>
                            Request to Join
                        </button>
                    {/if}
                </div>
            </div>
        </div>

        <div class="tabs">
            <button
                class="tab-button {activeTab === 'top' ? 'active' : ''}"
                onclick={() => handleTabChange("top")}
            >
                Top
            </button>
            <button
                class="tab-button {activeTab === 'latest' ? 'active' : ''}"
                onclick={() => handleTabChange("latest")}
            >
                Latest
            </button>
            <button
                class="tab-button {activeTab === 'media' ? 'active' : ''}"
                onclick={() => handleTabChange("media")}
            >
                Media
            </button>
            <button
                class="tab-button {activeTab === 'about' ? 'active' : ''}"
                onclick={() => handleTabChange("about")}
            >
                About
            </button>
            {#if community.role === "moderator" || community.role === "owner"}
                <button
                    class="tab-button {activeTab === 'manage' ? 'active' : ''}"
                    onclick={() => handleTabChange("manage")}
                >
                    Manage Members
                </button>
            {/if}
        </div>

        <div class="tab-content">
            {#if activeTab === "top"}
                <div class="top-content">
                    <div class="top-members">
                        <h3>Top Members</h3>
                        <div class="members-list">
                            {#each topMembers as member}
                                <!-- svelte-ignore a11y_click_events_have_key_events -->
                                <div
                                    role="button"
                                    tabindex="0"
                                    class="member-card"
                                    onclick={() =>
                                        navigateToProfile(member.username)}
                                >
                                    <img
                                        src={AVATAR_IMG +
                                            member.user_id +
                                            ".png" || "/placeholder.svg"}
                                        alt={member.username}
                                    />
                                    <span>{member.username}</span>
                                    <span class="followers">
                                        {member.followers ?? 0} followers
                                    </span>
                                </div>
                            {/each}
                        </div>
                    </div>

                    <div class="top-threads">
                        <h3>Most Liked Threads</h3>
                        {#each topThreads as thread}
                            <Post post={thread} />
                        {/each}
                    </div>
                </div>
            {:else if activeTab === "latest"}
                <div class="latest-threads">
                    {#each latestThreads as thread}
                        <Post post={thread} />
                    {/each}
                </div>
            {:else if activeTab === "media"}
                <div class="media-container">
                    {#if mediaThreads.length === 0}
                        <div class="empty-state">No media yet</div>
                    {:else}
                        <div class="media-grid">
                            {#each mediaThreads as mediaPost}
                                {#each mediaPost.media as media (media.media_url)}
                                    <!-- svelte-ignore a11y_click_events_have_key_events -->
                                    <!-- svelte-ignore a11y_no_static_element_interactions -->
                                    <div
                                        class="media-item"
                                        onclick={() =>
                                            navigateToThread(
                                                mediaPost.thread_id
                                            )}
                                    >
                                        {#if media.media_type === "image" || media.type === "gif"}
                                            <img
                                                src={`${THREAD_IMG}/${media.media_url}` ||
                                                    "/placeholder.svg"}
                                                alt="Media content"
                                            />
                                        {:else if media.type === "video"}
                                            <!-- svelte-ignore a11y_media_has_caption -->
                                            <video
                                                src={media.media_url}
                                                controls={false}
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
                                                    <polygon
                                                        points="5 3 19 12 5 21 5 3"
                                                    ></polygon>
                                                </svg>
                                            </div>
                                        {/if}
                                    </div>
                                {/each}
                            {/each}
                        </div>
                    {/if}
                </div>
            {:else if activeTab === "about"}
                <div class="about-content">
                    <div class="community-stats">
                        <h3>Community Information</h3>
                        <p>
                            <strong>Created At:</strong>
                            {community.created_at}
                        </p>
                        <p>
                            <strong>Created By:</strong>
                            {userProfiles[community.creator_id]
                                ? userProfiles[community.creator_id].username
                                : "Unknown"}
                        </p>
                        <p>
                            <strong>Member Count:</strong>
                            {community.member_count}
                        </p>
                    </div>

                    {#if community.rules}
                        <div class="community-rules">
                            <h3>Community Rules</h3>
                            <p>{community.rules}</p>
                        </div>
                    {/if}

                    <h3>Moderators</h3>
                    <div class="moderators-list">
                        {#each moderators as moderator}
                            <div class="member-card">
                                <!-- svelte-ignore a11y_click_events_have_key_events -->
                                <div
                                    class="moderator-info"
                                    role="button"
                                    tabindex="0"
                                    onclick={() =>
                                        navigateToProfile(
                                            userProfiles[moderator.user_id]
                                                .username
                                        )}
                                >
                                    <img
                                        src={AVATAR_IMG +
                                            moderator.user_id +
                                            ".png" || "/placeholder.svg"}
                                        alt={userProfiles[moderator.user_id]
                                            .username}
                                    />
                                    <span
                                        >{userProfiles[moderator.user_id]
                                            .username}</span
                                    >
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {:else if activeTab === "manage"}
                <div class="manage-content">
                    <div class="pending-requests">
                        <h3>Pending Join Requests</h3>
                        {#each pendingRequests as request}
                            <div class="member-card">
                                <img
                                    src={AVATAR_IMG +
                                        request.user_id +
                                        ".png" || "/placeholder.svg"}
                                    alt={userProfiles[request.user_id]
                                        ?.username}
                                    class="request-avatar"
                                />
                                <a
                                    href={"/profile/" +
                                        userProfiles[request.user_id]?.username}
                                >
                                    @{userProfiles[request.user_id]?.username}
                                </a>
                                <div class="actions">
                                    <button
                                        onclick={() =>
                                            handleAcceptRequest(
                                                request.user_id
                                            )}>Accept</button
                                    >
                                    <button
                                        onclick={() =>
                                            handleRejectRequest(
                                                request.user_id
                                            )}>Reject</button
                                    >
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/if}
        </div>
    </div>

    {#if showMembersModal}
        <MembersModal
            community_id={community.community_id}
            user_role={community.role}
            bind:is_open={showMembersModal}
        />
    {/if}
{:else}
    <div class="error">Community not found</div>
{/if}

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/community.scss";
    @use "../styles/profile.scss";
</style>
