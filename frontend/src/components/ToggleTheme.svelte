<script lang="ts">
    import { onMount } from "svelte";

    export let theme: "light" | "dark" = "dark";

    let currTheme: "light" | "dark" = theme;

    onMount(() => {
        const storedTheme = localStorage.getItem("theme") as "light" | "dark";

        if (storedTheme) {
            currTheme = storedTheme;
            theme = storedTheme;
            document.documentElement.setAttribute("data-theme", storedTheme);
        } else {
            localStorage.setItem("theme", theme);
            document.documentElement.setAttribute("data-theme", theme);
        }
    });

    function toggleTheme(): void {
        currTheme = currTheme === "light" ? "dark" : "light";
        localStorage.setItem("theme", currTheme);
        document.documentElement.setAttribute("data-theme", currTheme);
        theme = currTheme;
    }
</script>

<button class="theme-toggle" on:click={toggleTheme} aria-label="Toggle theme">
    {#if currTheme === "light"}
        <svg
            viewBox="0 0 24 24"
            width="24"
            height="24"
            stroke="currentColor"
            stroke-width="2"
            fill="none"
            stroke-linecap="round"
            stroke-linejoin="round"
        >
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
        </svg>
    {:else}
        <svg
            viewBox="0 0 24 24"
            width="24"
            height="24"
            stroke="currentColor"
            stroke-width="2"
            fill="none"
            stroke-linecap="round"
            stroke-linejoin="round"
        >
            <circle cx="12" cy="12" r="5"></circle>
            <line x1="12" y1="1" x2="12" y2="3"></line>
            <line x1="12" y1="21" x2="12" y2="23"></line>
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
            <line x1="1" y1="12" x2="3" y2="12"></line>
            <line x1="21" y1="12" x2="23" y2="12"></line>
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
        </svg>
    {/if}
</button>

<style>
    .theme-toggle {
        background: none;
        border: none;
        color: var(--text);
        cursor: pointer;
        padding: 8px;
        border-radius: 50%;
    }

    .theme-toggle:hover {
        background-color: rgba(128, 128, 128, 0.1);
    }
</style>
