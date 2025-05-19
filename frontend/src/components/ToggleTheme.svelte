<script lang="ts">
    import { onMount } from "svelte";
    import { Sun, Moon } from "@lucide/svelte";

    let currTheme: "light" | "dark" = $state("dark");

    onMount(() => {
        const storedTheme = localStorage.getItem("theme") as "light" | "dark";

        if (storedTheme) {
            currTheme = storedTheme;
            document.documentElement.setAttribute("data-theme", storedTheme);
        } else {
            localStorage.setItem("theme", currTheme);
            document.documentElement.setAttribute("data-theme", currTheme);
        }
    });

    function toggleTheme(): void {
        currTheme = currTheme === "light" ? "dark" : "light";
        localStorage.setItem("theme", currTheme);
        document.documentElement.setAttribute("data-theme", currTheme);
    }
</script>

<button class="theme-toggle" onclick={toggleTheme} aria-label="Toggle theme">
    {#if currTheme === "light"}
        <Moon />
    {:else}
        <Sun />
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
