<script lang="ts">
    import { onMount } from "svelte";
    import { toasts, removeToast } from "../stores/toast-wrapper";
    import Toast from "./Toast.svelte";
</script>

<div class="toast-container">
    {#each $toasts as toast (toast.id)}
        <Toast
            type={toast.type}
            title={toast.title}
            message={toast.message}
            duration={toast.duration}
            show={true}
        />
    {/each}
</div>

<style>
    .toast-container {
        position: fixed;
        top: 20px;
        right: 20px;
        z-index: 9999;
        display: flex;
        flex-direction: column;
        gap: 10px;
        max-width: 350px;
        width: 100%;
    }

    @media (max-width: 480px) {
        .toast-container {
            top: auto;
            bottom: 20px;
            left: 20px;
            right: 20px;
            max-width: none;
        }
    }
</style>
