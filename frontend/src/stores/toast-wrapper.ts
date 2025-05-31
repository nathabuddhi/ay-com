import { writable } from "svelte/store";

export type ToastType = "error" | "success" | "info";

export interface Toast {
    id: string;
    type: ToastType;
    title: string;
    message: string;
    duration: number;
}

export const toasts = writable<Toast[]>([]);

export function addToast(
    type: ToastType,
    message: string,
    title: string = "",
    duration: number = 5000
): string {
    const id = generateId();
    const toast: Toast = {
        id,
        type,
        title,
        message,
        duration,
    };

    toasts.update((all) => [toast, ...all]);

    if (duration > 0) {
        setTimeout(() => {
            removeToast(id);
        }, duration);
    }

    return id;
}

export function removeToast(id: string): void {
    setTimeout(() => {
        toasts.update((all) => all.filter((t) => t.id !== id));
    }, 300);
}

export function clearToasts(): void {
    toasts.set([]);
}

function generateId(): string {
    return Math.random().toString(36).substring(2, 9);
}
