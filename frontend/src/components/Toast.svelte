<script lang="ts">
    import { onMount } from "svelte";

    const { type, title, message, duration, show } = $props<{
        type: "error" | "success" | "info";
        title: string;
        message: string;
        duration: number;
        show: boolean;
    }>();

    let isVisible = $state(false);
    let timeoutId: number;

    $effect(() => {
        if (show) {
            showToast();
        }
    });

    function showToast(): void {
        isVisible = true;

        if (duration > 0) {
            clearTimeout(timeoutId);
            timeoutId = setTimeout(() => {
                closeToast();
            }, duration) as number;
        }
    }

    function closeToast(): void {
        isVisible = false;
        clearTimeout(timeoutId);
    }

    onMount(() => {
        return () => {
            clearTimeout(timeoutId);
        };
    });
</script>

<div
    class="toast {type} {isVisible ? 'show' : 'hide'}"
    role="alert"
    aria-live="assertive"
>
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
