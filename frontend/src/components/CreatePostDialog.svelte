<script lang="ts">
    import { X } from "@lucide/svelte";
    import CreatePostForm from "./CreatePostForm.svelte";

    let { isOpen = $bindable() }: { isOpen: boolean } = $props();
</script>

{#if isOpen}
    <div class="create-post-background">
        <div class="create-post">
            <button class="close-button" onclick={() => (isOpen = false)}>
                <X />
            </button>
            <CreatePostForm {isOpen} mode={"post"} />
        </div>
    </div>
{/if}

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/create-post.scss";
</style>
