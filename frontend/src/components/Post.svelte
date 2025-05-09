<script lang="ts">
    import { formatDistanceToNow } from "date-fns";
    import {
        BadgeCheck,
        Share,
        MessageCircleMore,
        Repeat2,
        Heart,
        EllipsisVertical,
    } from "@lucide/svelte";
    import type { Thread } from "../types/thread";
    import type { UserProfile } from "../types/user";

    export let post: Thread;
    let user: UserProfile = {
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
    };

    const formattedTime = formatDistanceToNow(new Date(post.timestamp), {
        addSuffix: true,
    });

    function formatNumber(num: number): string {
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + "M";
        } else if (num >= 1000) {
            return (num / 1000).toFixed(1) + "K";
        } else {
            return num.toString();
        }
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

    function handleMentionClick(event: MouseEvent, username: string) {
        event.preventDefault();
        event.stopPropagation();
        window.location.href = `/profile/${username}`;
    }

    function handleHashtagClick(event: MouseEvent, hashtag: string) {
        event.preventDefault();
        event.stopPropagation();
        window.location.href = `/explore?q=%23${hashtag}`;
    }

    function addEventListeners(node: HTMLElement) {
        const mentions = node.querySelectorAll(".mention");
        mentions.forEach((mention) => {
            mention.addEventListener("click", (event) => {
                const username = mention.textContent?.substring(1);
                if (username) handleMentionClick(event as MouseEvent, username);
            });
        });

        const hashtags = node.querySelectorAll(".hashtag");
        hashtags.forEach((hashtag) => {
            hashtag.addEventListener("click", (event) => {
                const tag = hashtag.textContent?.substring(1);
                if (tag) handleHashtagClick(event as MouseEvent, tag);
            });
        });

        return {
            destroy() {},
        };
    }
</script>

<article class="post">
    <div class="post-avatar">
        <img
            src={`http://localhost:5000/images/profile/${post.user_id}`}
            alt={post.user_id}
        />
    </div>

    <div class="post-content">
        <div class="post-header">
            <div class="post-user-info">
                <span class="post-user-name">
                    {user.username}
                    {#if user.is_verified}
                        <BadgeCheck />
                    {/if}
                </span>
                <span class="post-user-username">@{user.username}</span>
                <span class="post-time">{formattedTime}</span>
            </div>

            <button class="post-more-options">
                <EllipsisVertical />
            </button>
        </div>

        <div class="post-text" use:addEventListeners>
            {@html processContent(post.content)}
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
                            src={`${import.meta.env.MEDIA_LINK}/{image.media_url}`}
                            alt="Post image {i + 1}"
                        />
                    </div>
                {/each}
            </div>
        {/if}

        <div class="post-actions">
            <button class="post-action comment">
                <MessageCircleMore />
                <span>{formatNumber(post.reply_count)}</span>
            </button>

            <button class="post-action repost">
                <Repeat2 />
                <span>{formatNumber(post.repost_count)}</span>
            </button>

            <button class="post-action like">
                <Heart />
                <span>{formatNumber(post.like_count)}</span>
            </button>

            <button class="post-action share">
                <Share />
            </button>
        </div>
    </div>
</article>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
