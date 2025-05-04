import { addToast } from "../stores/toast-wrapper";
import type { BoolPayload } from "../types/api";

export async function isLoggedIn(): Promise<boolean> {
    if (!getToken()) return false;

    return await checkTokenValidity();
}

export function getToken(): string {
    return localStorage.getItem("token") || "";
}

export function setToken(token: string): void {
    localStorage.setItem("token", token);
}

export function removeToken(): void {
    localStorage.removeItem("token");
}

export function logout(): void {
    removeToken();
    window.location.href = "/";
}

async function checkTokenValidity(): Promise<boolean> {
    addToast("info", "Checking user cookie validity...", "Token Check");

    const token = getToken();
    if (!token) return false;
    const response = await fetch("http://localhost:5000/user/checktoken", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: token,
        },
        body: JSON.stringify({}),
    });
    const data: BoolPayload = await response.json();
    if (!response.ok || !data.value) {
        removeToken();
        return false;
    }

    return true;
}
