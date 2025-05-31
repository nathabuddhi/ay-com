<script lang="ts">
    const { type, title, message, id } = $props<{
        type: "error" | "success" | "info";
        title: string;
        message: string;
        id: string;
    }>();

    import { removeToast } from "../stores/toast-wrapper";

    function closeToast(): void {
        removeToast(id);
    }
</script>

<div class="toast {type}" role="alert" aria-live="assertive">
    <div class="toast-content">
        {#if title}
            <div class="toast-title">{title}</div>
        {/if}
        <div class="toast-message">{message}</div>
    </div>
    <button
        class="toast-close"
        onclick={closeToast}
        aria-label="Close notification"
    >
        &times;
    </button>
</div>

<!-- svelte-ignore css_unused_selector -->
<style lang="scss">
    @use "../styles/toast.scss";
</style>
