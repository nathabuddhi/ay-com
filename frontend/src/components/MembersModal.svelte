<script lang="ts">
    import { onMount } from "svelte";
    import type { Community, CommunityMember } from "../types/community";
    import {
        demoteModerator,
        getCommunityMembers,
        promoteMember,
    } from "../controllers/community-controller";
    import type { UserProfile } from "../types/user";
    import { addToast } from "../stores/toast-wrapper";
    import { getUserById } from "../controllers/user-controller";
    import { AVATAR_IMG } from "../env_var";
    import { X } from "@lucide/svelte";
    import Pagination from "./Pagination.svelte";

    let {
        community_id,
        user_role = $bindable<"member" | "moderator" | "owner">("member"),
        is_open = $bindable(),
    }: { community_id: string; user_role: string; is_open: boolean } = $props();

    let activeTab = $state<"members" | "moderators">("members");
    let members = $state<CommunityMember[]>([]);
    let moderators = $state<CommunityMember[]>([]);
    let userProfiles = $state<Record<string, UserProfile>>({});
    let searchQuery = $state("");

    onMount(async () => {
        await loadMembers();
    });

    async function loadMembers() {
        try {
            const response = await getCommunityMembers(community_id);

            if (
                response.success &&
                response.payload &&
                response.payload.members
            ) {
                members = response.payload.members;

                const profilesArr = await Promise.all(
                    members.map(async (member) => {
                        const user = await getUserById(member.user_id);
                        return { user_id: member.user_id, user };
                    })
                );
                userProfiles = profilesArr.reduce(
                    (acc, curr) => {
                        if (curr.user) {
                            acc[curr.user_id] = curr.user;
                        }
                        return acc;
                    },
                    {} as Record<string, UserProfile>
                );

                moderators = members.filter(
                    (member) => member.role === "moderator"
                );
                members = members.filter((member) => member.role === "member");
            }
        } catch (error) {
            addToast("error", "Failed to load members", "Error");
            members = [];
        }
    }

    function handleUserClick(user: UserProfile) {
        if (user) {
            window.location.href = `/profile/${user.username}`;
        }
    }

    async function handlePromoteMember(user_id: string) {
        const response = await promoteMember(community_id, user_id);
        if (response.success) {
            addToast("success", "Member promoted to moderator!", "Success");
            await loadMembers();
        } else {
            addToast("error", response.message, "Promotion Failed");
        }
    }

    async function handleDemoteMember(user_id: string) {
        const response = await demoteModerator(community_id, user_id);
        if (response.success) {
            addToast("success", "Moderator demoted to member!", "Success");
            await loadMembers();
        } else {
            addToast("error", response.message, "Demotion Failed");
        }
    }
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="modal-overlay">
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
        <div class="modal-header">
            <h2>Community Members</h2>
            <button class="close-btn" onclick={() => (is_open = false)}>
                <X />
            </button>
        </div>

        <div class="search-bar">
            <input
                type="text"
                placeholder="Search members..."
                bind:value={searchQuery}
            />
        </div>
        <div class="tabs">
            <button
                class="tab-button {activeTab === 'members' ? 'active' : ''}"
                onclick={() => (activeTab = "members")}
            >
                Members
            </button>
            <button
                class="tab-button {activeTab === 'moderators' ? 'active' : ''}"
                onclick={() => (activeTab = "moderators")}
            >
                Moderators
            </button>
        </div>
        <div class="members-list">
            {#if activeTab === "moderators"}
                {#each moderators.filter((m) => !searchQuery || userProfiles[m.user_id]?.username
                            .toLowerCase()
                            .includes(searchQuery.toLowerCase())) as member}
                    {#if userProfiles[member.user_id]}
                        <div class="member-item">
                            <div
                                class="member-info"
                                role="button"
                                tabindex="0"
                                onclick={() =>
                                    handleUserClick(
                                        userProfiles[member.user_id]
                                    )}
                            >
                                <img
                                    src={AVATAR_IMG + member.user_id + ".png" ||
                                        "/placeholder.svg"}
                                    alt={userProfiles[member.user_id].username}
                                />
                                <span
                                    >{userProfiles[member.user_id]
                                        .username}</span
                                >
                            </div>
                            {#if user_role === "owner"}
                                <button
                                    onclick={() =>
                                        handleDemoteMember(member.user_id)}
                                >
                                    Demote to Member
                                </button>
                            {/if}
                        </div>
                    {/if}
                {/each}
                <Pagination totalItems={moderators.length} />
            {:else}
                {#each members.filter((m) => !searchQuery || userProfiles[m.user_id]?.username
                            .toLowerCase()
                            .includes(searchQuery.toLowerCase())) as member}
                    {#if userProfiles[member.user_id]}
                        <div class="member-item">
                            <div
                                class="member-info"
                                role="button"
                                tabindex="0"
                                onclick={() =>
                                    handleUserClick(
                                        userProfiles[member.user_id]
                                    )}
                            >
                                <img
                                    src={AVATAR_IMG + member.user_id + ".png" ||
                                        "/placeholder.svg"}
                                    alt={userProfiles[member.user_id].username}
                                />
                                <span
                                    >{userProfiles[member.user_id]
                                        .username}</span
                                >
                            </div>
                            {#if user_role === "owner"}
                                <button
                                    onclick={() =>
                                        handlePromoteMember(member.user_id)}
                                >
                                    Promote to Moderator
                                </button>
                            {/if}
                        </div>
                    {/if}
                {/each}
                <Pagination totalItems={members.length} />
            {/if}
        </div>
    </div>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/community.scss";
</style>
