<script lang="ts">
    import { formatDistanceToNow } from "date-fns";
    import {
        BadgeCheck,
        Share,
        MessageCircleMore,
        Repeat2,
        Heart,
        EllipsisVertical,
        Bookmark,
    } from "@lucide/svelte";
    import type { Thread } from "../types/thread";
    import type { UserProfile } from "../types/user";
    import { onMount } from "svelte";
    import {
        getThreadOwner,
        toggleBookmark,
        toggleLike,
        toggleRepost,
    } from "../controllers/thread-controller";
    import { AVATAR_IMG, THREAD_IMG } from "../env_var";
    import { addToast } from "../stores/toast-wrapper";

    let { post } = $props<{ post: Thread }>();

    async function handleLikeClick() {
        if (post.is_liking) {
            const response = await toggleLike(post.thread_id);
            if (response.success) {
                post.is_liking = false;
                post.like_count = post.like_count ? post.like_count - 1 : 0;
            } else {
                addToast(
                    "error",
                    "Failed unliking post: " + response.message,
                    "Error!"
                );
            }
        } else {
            const response = await toggleLike(post.thread_id);
            if (response.success) {
                post.is_liking = true;
                post.like_count = post.like_count ? post.like_count + 1 : 1;
            } else {
                addToast(
                    "error",
                    "Failed liking post: " + response.message,
                    "Error!"
                );
            }
        }
    }

    async function handleRepostClick() {
        if (post.is_reposting) {
            const response = await toggleRepost(post.thread_id);
            if (response.success) {
                post.is_reposting = false;
                post.repost_count = post.repost_count
                    ? post.repost_count - 1
                    : 0;
            } else {
                addToast(
                    "error",
                    "Failed unreposting post: " + response.message,
                    "Error!"
                );
            }
        } else {
            const response = await toggleRepost(post.thread_id);
            if (response.success) {
                post.is_reposting = true;
                post.repost_count = post.repost_count
                    ? post.repost_count + 1
                    : 1;
            } else {
                addToast(
                    "error",
                    "Failed reposting post: " + response.message,
                    "Error!"
                );
            }
        }
    }

    async function handleBookmarkClick() {
        if (post.is_bookmarking) {
            const response = await toggleBookmark(post.thread_id);
            if (response.success) {
                post.is_bookmarking = false;
            } else {
                addToast(
                    "error",
                    "Failed unliking post: " + response.message,
                    "Error!"
                );
            }
        } else {
            const response = await toggleBookmark(post.thread_id);
            if (response.success) {
                post.is_bookmarking = true;
            } else {
                addToast(
                    "error",
                    "Failed liking post: " + response.message,
                    "Error!"
                );
            }
        }
    }

    let user = $state<UserProfile>({
        name: "loading",
        username: "loading",
        is_verified: false,
        user_id: post.user_id,
        bio: "",
        followers: 0,
        following: 0,
        gender: "",
        date_of_birth: "",
        email: "",
        join_date: "",
    });

    onMount(async () => {
        const response = await getThreadOwner(post.user_id);
        console.log("Fetching thread owner with user_id:", post.user_id);
        if (response) {
            user = response;
        } else {
            user = {
                name: "Unknown.",
                username: "Unknown.",
                is_verified: false,
                user_id: post.user_id,
                bio: "",
                followers: 0,
                following: 0,
                gender: "",
                date_of_birth: "",
                email: "",
                join_date: "",
            };
        }
    });

    function formatNumber(num: number): string {
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + "M";
        } else if (num >= 1000) {
            return (num / 1000).toFixed(1) + "K";
        } else return String(num);
    }

    function processContent(content: string): string {
        let processed = content.replace(
            /@(\w+)/g,
            '<span class="mention">@$1</span>'
        );

        processed = processed.replace(
            /#(\w+)/g,
            '<span class="hashtag">#$1</span>'
        );

        return processed;
    }

    function handleMentionClick(event: MouseEvent) {
        if (user.username == "loading") return;
        event.preventDefault();
        event.stopPropagation();
        window.location.href = `/profile/${user.username}`;
    }

    function handleHashtagClick(hashtag: string) {
        window.location.href = `/explore?q=%23${hashtag}`;
    }
</script>

<article class="post">
    <div class="post-avatar">
        <img src={`${AVATAR_IMG}/${post.user_id}.png`} alt={post.user_id} />
    </div>

    <div class="post-content">
        <div class="post-header">
            <div class="post-user-info">
                <div>
                    <span class="post-user-name">
                        {user.name}
                        {#if user.is_verified}
                            <BadgeCheck />
                        {/if}
                    </span>
                    <button
                        class="post-user-username"
                        onclick={handleMentionClick}>@{user.username}</button
                    >
                </div>
                <span class="post-time">{post.posted_at}</span>
            </div>

            <button class="post-more-options">
                <EllipsisVertical />
            </button>
        </div>

        {#if post.media.length > 0}
            <div
                class="post-images {post.media.length > 1
                    ? 'multiple-images'
                    : ''}"
            >
                {#each post.media as image, i}
                    <div class="image-container">
                        <img
                            src={`${THREAD_IMG}/${image.media_url}`}
                            alt="Post image {i + 1}"
                        />
                    </div>
                {/each}
            </div>
        {/if}

        <div class="post-actions">
            <button class="post-action like" onclick={handleLikeClick}>
                <Heart fill={post.is_liking ? "red" : ""} />
                <span>{formatNumber(post.like_count ?? 0)}</span>
            </button>

            <button class="post-action comment">
                <MessageCircleMore />
                <span>{formatNumber(post.reply_count ?? 0)}</span>
            </button>

            <button class="post-action repost" onclick={handleRepostClick}>
                <Repeat2 fill={post.is_reposting ? "green" : ""} />
                <span>{formatNumber(post.repost_count ?? 0)}</span>
            </button>

            <button class="post-action bookmark" onclick={handleBookmarkClick}>
                <Bookmark fill={post.is_bookmarking ? "gold" : ""} />
            </button>
        </div>
    </div>
</article>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
