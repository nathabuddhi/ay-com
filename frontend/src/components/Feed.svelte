<script lang="ts">
    import Post from "./Post.svelte";
    // import CreatePost from "./CreatePost.svelte";

    // Dummy data for posts
    const posts = [
        {
            id: "1",
            user: {
                id: "101",
                name: "Deedee",
                username: "deedeesmith",
                is_verified: true,
            },
            content:
                "I've been working on my home office setup! What do you think? #homeoffice #workfromhome",
            images: ["/images/post1.jpg", "/images/post2.jpg"],
            timestamp: "2023-05-04T14:30:00Z",
            likes: 235,
            reposts: 36,
            comments: 128,
        },
        {
            id: "2",
            user: {
                id: "102",
                name: "BAZNAS",
                username: "baznasindonesia",
                is_verified: true,
            },
            content:
                "Mari tunjukan kepedulian untuk membantu masyarakat Indonesia #BantuanSosial #Indonesia",
            images: [],
            timestamp: "2023-05-04T12:15:00Z",
            likes: 456,
            reposts: 89,
            comments: 32,
        },
        {
            id: "3",
            user: {
                id: "103",
                name: "Tech Insider",
                username: "techinsider",
                is_verified: false,
            },
            content:
                "New features coming to @AYplatform next week! Stay tuned for the big announcement. #tech #innovation",
            images: ["/images/post3.jpg"],
            timestamp: "2023-05-04T10:45:00Z",
            likes: 789,
            reposts: 145,
            comments: 67,
        },
    ];

    const tabs = ["For you", "Following"];
    let activeTab = "For you";

    function setActiveTab(tab: string) {
        activeTab = tab;
    }
</script>

<div class="feed-container">
    <div class="feed-header">
        <h1>Home</h1>
        <div class="feed-tabs">
            {#each tabs as tab}
                <button
                    class="tab-button {activeTab === tab ? 'active' : ''}"
                    on:click={() => setActiveTab(tab)}
                >
                    {tab}
                </button>
            {/each}
        </div>
    </div>

    <div class="posts-container">
        {#each posts as post (post.id)}
            <Post {post} />
        {/each}
    </div>
</div>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
    @use "../styles/home.scss";
</style>
