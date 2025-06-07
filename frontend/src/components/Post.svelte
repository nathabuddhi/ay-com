<script lang="ts">
    import {
        BadgeCheck,
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
        deleteThread,
        pinThread,
        toggleBookmark,
        toggleLike,
        toggleRepost,
        voteThreadPoll,
    } from "../controllers/thread-controller";
    import { AVATAR_IMG, THREAD_IMG } from "../env_var";
    import { addToast } from "../stores/toast-wrapper";
    import { processContent } from "../controllers/util";
    import { getUserById } from "../controllers/user-controller";

    let { post }: { post: Thread } = $props<{ post: Thread }>();
    let isLoading: boolean = $state<boolean>(true);
    let isMoreOptionsOpen: boolean = $state<boolean>(false);

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
        is_private: false,
    });

    onMount(async () => {
        const response = await getUserById(post.user_id);
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
                is_private: false,
            };
        }

        isLoading = false;
    });

    function togglePopover() {
        isMoreOptionsOpen = !isMoreOptionsOpen;
    }

    async function handleDeleteClick() {
        const response = await deleteThread(post.thread_id);

        if (response.success) {
            addToast("success", "Post deleted successfully!", "Success!");
            setTimeout(() => {
                window.location.reload();
            }, 500);
        } else {
            addToast(
                "error",
                "Failed voting on poll: " + response.message,
                "Error!"
            );
        }
    }

    async function handlePinClick() {
        const response = await pinThread(post.thread_id);

        if (response.success) {
            addToast("success", "Post pinned successfully!", "Success!");
            setTimeout(() => {
                window.location.reload();
            }, 500);
        } else {
            addToast(
                "error",
                "Failed pinning thread: " + response.message,
                "Error!"
            );
        }
    }

    async function handleShareClick() {}

    async function handleVoteClick(poll: string) {
        const response = await voteThreadPoll(post.thread_id, poll);

        if (response.success) {
            post.poll_options = post.poll_options.map((option) => {
                if (option.option === poll) {
                    if (!option.is_voting)
                        return {
                            option: option.option,
                            vote_count: option.vote_count
                                ? option.vote_count + 1
                                : 1,
                            is_voting: true,
                        };
                    else
                        return {
                            option: option.option,
                            vote_count: (option.vote_count ?? 0) - 1,
                            is_voting: false,
                        };
                }
                return {
                    option: option.option,
                    vote_count: option.is_voting
                        ? option.vote_count - 1
                        : option.vote_count,
                    is_voting: false,
                };
            });
        } else {
            addToast(
                "error",
                "Failed voting on poll: " + response.message,
                "Error!"
            );
        }
    }

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

    function redirectToThreadDetail() {
        if (window.location.pathname.includes("thread")) return;

        window.location.href = `/thread/${post.thread_id}`;
    }

    function formatNumber(num: number): string {
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + "M";
        } else if (num >= 1000) {
            return (num / 1000).toFixed(1) + "K";
        } else return String(num);
    }

    function handleMentionClick(event: MouseEvent) {
        if (user.username == "loading") return;
        event.preventDefault();
        event.stopPropagation();
        window.location.href = `/profile/${user.username}`;
    }
</script>

<p class="repost-text">{post.is_reposting && "You Reposted"}</p>
{#if isLoading}
    <article class="post loading-skeleton">
        <div class="post-avatar">
            <div class="skeleton skeleton-avatar"></div>
        </div>

        <div class="post-content">
            <div class="post-header">
                <div class="post-user-info">
                    <div>
                        <div class="skeleton skeleton-name"></div>
                        <div class="skeleton skeleton-username"></div>
                    </div>
                    <div>
                        <div class="skeleton skeleton-category"></div>
                        <div class="skeleton skeleton-timestamp"></div>
                    </div>
                </div>
                <div class="post-more-options-container">
                    <div class="skeleton skeleton-ellipsis"></div>
                </div>
            </div>

            <div class="post-text">
                <div class="skeleton skeleton-text-line"></div>
                <div class="skeleton skeleton-text-line short"></div>
            </div>

            <div class="post-images multiple-images">
                <div class="skeleton skeleton-image"></div>
                <div class="skeleton skeleton-image"></div>
            </div>

            <div class="post-polls">
                <div class="skeleton skeleton-poll"></div>
                <div class="skeleton skeleton-poll"></div>
            </div>

            <div class="post-actions">
                <div class="skeleton skeleton-action"></div>
                <div class="skeleton skeleton-action"></div>
                <div class="skeleton skeleton-action"></div>
                <div class="skeleton skeleton-action"></div>
            </div>
        </div>
    </article>
{:else}
    <article class="post">
        <div class="post-avatar">
            <img src={`${AVATAR_IMG}/${post.user_id}.png`} alt={post.user_id} />
        </div>

        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
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
                            onclick={handleMentionClick}
                            >@{user.username}</button
                        >
                    </div>
                    <div>
                        <span class="post-time">Category: {post.category}</span>
                        <br />
                        <span class="post-time">{post.posted_at}</span>
                    </div>
                </div>

                <div class="post-more-options-container">
                    <button class="post-more-options" onclick={togglePopover}>
                        <EllipsisVertical />
                    </button>

                    {#if isMoreOptionsOpen}
                        <div class="popover" role="menu">
                            {#if post.user_id === localStorage.getItem("user_id")}
                                <button onclick={handleDeleteClick}>
                                    Delete
                                </button>
                                <button onclick={handlePinClick}>
                                    {post.pinned ? "Pin" : "Unpin"}
                                </button>
                            {/if}
                            <button onclick={handleShareClick}>Share</button>
                        </div>
                    {/if}
                </div>
            </div>
            <div class="post-text">
                {@html processContent(post.content)}
            </div>
            {#if post.media && post.media.length > 0}
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
            {#if post.poll_options && post.poll_options.length > 0}
                <div class="post-polls">
                    {#each post.poll_options as poll, i}
                        <button
                            class={`poll-container ${poll.is_voting ? "active" : ""}`}
                            onclick={() => handleVoteClick(poll.option)}
                        >
                            <span class={`poll-option`}>{poll.option}</span>
                            <span class="poll-votes"
                                >{poll.vote_count ?? 0} votes</span
                            >
                        </button>
                    {/each}
                </div>
            {/if}

            <div class="post-actions">
                <button class="post-action like" onclick={handleLikeClick}>
                    <Heart fill={post.is_liking ? "red" : ""} />
                    <span>{formatNumber(post.like_count ?? 0)}</span>
                </button>

                <button
                    class="post-action comment"
                    onclick={redirectToThreadDetail}
                >
                    <MessageCircleMore />
                    <span>{formatNumber(post.reply_count ?? 0)}</span>
                </button>

                <button class="post-action repost" onclick={handleRepostClick}>
                    <Repeat2 fill={post.is_reposting ? "green" : ""} />
                    <span>{formatNumber(post.repost_count ?? 0)}</span>
                </button>

                <button
                    class="post-action bookmark"
                    onclick={handleBookmarkClick}
                >
                    <Bookmark fill={post.is_bookmarking ? "gold" : ""} />
                </button>
            </div>
        </div>
    </article>
{/if}

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
