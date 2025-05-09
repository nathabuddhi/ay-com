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
    localStorage.removeItem("token");
    localStorage.removeItem("name");
    localStorage.removeItem("username");
    localStorage.removeItem("is_verified");
    localStorage.removeItem("user_id");

    window.location.href = "/";
}

async function checkTokenValidity(): Promise<boolean> {
    if (
        localStorage.getItem("token_checked") === "true" &&
        localStorage.getItem("token_checked_expiry")
    ) {
        const expiryTime = parseInt(
            localStorage.getItem("token_checked_expiry") || "0"
        );
        if (Date.now() < expiryTime) {
            return true;
        } else {
            localStorage.removeItem("token_checked");
            localStorage.removeItem("token_checked_expiry");
        }
        return true;
    }

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
    } else {
        const expiryTime = Date.now() + 20 * 60 * 1000;
        localStorage.setItem("token_checked", "true");
        localStorage.setItem("token_checked_expiry", expiryTime.toString());
    }

    return true;
}
