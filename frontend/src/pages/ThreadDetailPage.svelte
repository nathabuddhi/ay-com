<script lang="ts">
    import { onMount } from "svelte";
    import type { Thread, ThreadReply } from "../types/thread";
    import Post from "../components/Post.svelte";
    import {
        getThreadById,
        replyToThread,
    } from "../controllers/thread-controller";
    import { addToast } from "../stores/toast-wrapper";
    import { AVATAR_IMG } from "../env_var";
    import type { UserProfile } from "../types/user";
    import { getUserById } from "../controllers/user-controller";
    import { BadgeCheck } from "@lucide/svelte";

    let post: Thread = $state<Thread>({
        thread_id: "",
        user_id: "",
        content: "",
        category: "",
        is_advertisement: false,
        reply_permission: "",
        pinned: false,
        media: [],
        poll_options: [],
        like_count: 0,
        repost_count: 0,
        reply_count: 0,
        posted_at: "",
        is_private: false,
        is_liking: false,
        is_bookmarking: false,
        is_reposting: false,
    });
    let replies: ThreadReply[] = $state([]);
    let userProfiles: Record<string, UserProfile> = $state({});

    let newReply: string = $state<string>("");

    onMount(async () => {
        const thread_id = window.location.pathname.split("/").pop() ?? "";

        const response = await getThreadById(thread_id);

        if (!response.success) {
            addToast("error", response.message, "Error fetching tread!");
            window.location.href = "/home";
        } else if (response.payload) {
            post = response?.payload.thread;
            replies = response.payload.replies;

            const userIds = Array.from(
                new Set(response.payload.replies.map((r) => r.user_id))
            );
            const responses = await Promise.all(
                userIds.map((id) => getUserById(id))
            );

            const newProfiles: Record<string, UserProfile> = {};
            userIds.forEach((id, index) => {
                const res = responses[index];
                if (res) {
                    newProfiles[id] = res;
                }
            });

            userProfiles = newProfiles;
        }
    });

    async function handleReplyClick() {
        const response = await replyToThread(post.thread_id, newReply);

        if (!response.success) {
            addToast("error", response.message, "Error replying to thread!");
            return;
        }
        addToast("success", "Reply added successfully!", "Reply Success");
        window.location.reload();
    }

    async function handlePinClick(reply: ThreadReply) {
        const response = await pinReply(replyId);

        if (!response.success) {
            addToast("error", response.message, "Error pinning reply!");
            return;
        }
        addToast("success", "Reply pinned successfully!", "Pin Success");
        window.location.reload();
    }

    async function handleDeleteClick(replyId: string) {
        const response = await deleteReply(replyId);

        if (!response.success) {
            addToast("error", response.message, "Error deleting reply!");
            return;
        }
        addToast("success", "Reply deleted successfully!", "Delete Success");
        window.location.reload();
    }
</script>

{#if post.thread_id === ""}
    <div class="loading">Loading thread...</div>
{:else}
    <Post {post} />
    <div class="create-reply">
        <input type="text" placeholder="Write a Reply" bind:value={newReply} />
        <button class="reply-button" onclick={handleReplyClick}>Reply</button>
    </div>
    <div class="replies-list">
        {#each replies as reply}
            <div class="reply-item">
                <img
                    src={`${AVATAR_IMG}/${reply?.user_id}.png`}
                    alt="Avatar"
                    class="avatar"
                />
                <div class="reply-content">
                    <div class="reply-header">
                        <span class="reply-name"
                            >{userProfiles[reply.user_id]?.name}</span
                        >
                        {#if userProfiles[reply.user_id]?.is_verified}
                            <BadgeCheck />
                        {/if}
                        <span class="reply-username"
                            >@{userProfiles[reply.user_id]?.username}</span
                        >
                        <span class="reply-time">{reply.timestamp}</span>
                    </div>
                    <div class="reply-text-container">
                        <div class="reply-text">{reply.content}</div>
                        <div>
                            {#if localStorage.getItem("user_id") === post.user_id}
                                <button
                                    class="pin-button"
                                    onclick={() => {
                                        handlePinClick(reply);
                                    }}
                                    >{reply.is_pinned ? "Unpin" : "Pin"}</button
                                >
                            {/if}
                            {#if localStorage.getItem("user_id") === post.user_id || localStorage.getItem("user_id") === reply.user_id}
                                <button
                                    class="delete-button"
                                    onclick={() => {
                                        handleDeleteClick(reply.id);
                                    }}>Delete</button
                                >{/if}
                        </div>
                    </div>
                </div>
            </div>
        {/each}
    </div>
{/if}

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
