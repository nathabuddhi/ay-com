<script lang="ts">
    import { onMount } from "svelte";
    import type { Thread } from "../types/thread";
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
    import CreatePostForm from "../components/CreatePostForm.svelte";

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
        community_id: "",
    });
    let replies: Thread[] = $state([]);
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
</script>

{#if post.thread_id === ""}
    <div class="loading">Loading thread...</div>
{:else}
    <Post {post} />
    <CreatePostForm mode="comment" threadId={post.thread_id} isOpen={true} />
    <div class="replies-list">
        {#each replies as post}
            <Post {post} />
        {/each}
    </div>
{/if}

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
