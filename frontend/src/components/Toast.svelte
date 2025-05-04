<script lang="ts">
    import { onMount, createEventDispatcher } from "svelte";

    export let type: "error" | "success" | "info" = "info";
    export let title: string = "";
    export let message: string = "";
    export let duration: number = 5000;
    export let show: boolean = false;

    let isVisible = false;
    let timeoutId: number;
    const dispatch = createEventDispatcher<{ close: void }>();

    $: if (show) {
        showToast();
    }

    function showToast(): void {
        isVisible = true;

        if (duration > 0) {
            clearTimeout(timeoutId);
            timeoutId = setTimeout(() => {
                closeToast();
            }, duration) as unknown as number;
        }
    }

    function closeToast(): void {
        isVisible = false;
        clearTimeout(timeoutId);

        setTimeout(() => {
            dispatch("close");
        }, 300);
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
        on:click={closeToast}
        aria-label="Close notification"
    >
        &times;
    </button>
</div>

<style lang="scss">
    @use "../styles/toast.scss";
</style>
