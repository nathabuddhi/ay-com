import { writable } from "svelte/store";

type Theme = "light" | "dark";

export const theme = writable<Theme>("dark");

export function toggleTheme(): void {
    theme.update((currentTheme) => {
        const newTheme = currentTheme === "light" ? "dark" : "light";
        if (typeof localStorage !== "undefined") {
            localStorage.setItem("theme", newTheme);
        }
        if (typeof document !== "undefined") {
            document.documentElement.setAttribute("data-theme", newTheme);
        }
        return newTheme;
    });
}

export function initTheme(): void {
    if (typeof window !== "undefined") {
        const savedTheme = (localStorage.getItem("theme") as Theme) || "dark";
        theme.set(savedTheme);
        document.documentElement.setAttribute("data-theme", savedTheme);
    }
}
