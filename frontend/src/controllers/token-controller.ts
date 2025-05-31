import { addToast } from "../stores/toast-wrapper";
import type { ApiResponse } from "../types/api";
import type { LoginResponse } from "../types/user";

function setCookie(name: string, value: string, minutes: number): void {
    const expires = new Date(Date.now() + minutes * 60 * 1000).toUTCString();
    document.cookie = `${name}=${value}; expires=${expires}; path=/; SameSite=Strict`;
}

function getCookie(name: string): string {
    const cookies = document.cookie.split("; ");
    for (const cookie of cookies) {
        const [key, val] = cookie.split("=");
        if (key === name) return decodeURIComponent(val);
    }
    return "";
}

function deleteCookie(name: string): void {
    document.cookie = `${name}=; Max-Age=0; path=/; SameSite=Strict`;
}

export async function getToken(): Promise<string> {
    const token = getCookie("token");
    return token;
}

export function setToken(token: string): void {
    setCookie("token", token, 20);
}

export function setRefreshToken(token: string): void {
    setCookie("refresh_token", token, 60 * 24 * 7);
}

export function removeToken(): void {
    deleteCookie("token");
}

export async function isLoggedIn(): Promise<boolean> {
    return await checkTokenValidity();
}

export function logout(): void {
    ["token", "refresh_token"].forEach(deleteCookie);
    localStorage.removeItem("user_id");
    localStorage.removeItem("username");
    localStorage.removeItem("name");
    localStorage.removeItem("is_verified");
    window.location.href = "/";
}

async function refreshToken(refresh_token: string): Promise<string> {
    const response = await fetch("http://localhost:5000/user/refreshtoken", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ value: refresh_token }),
    });

    const data: ApiResponse<LoginResponse> = await response.json();
    if (!response.ok || !data.success) {
        logout();
        return "";
    } else {
        if (data.payload?.token) {
            setToken(data.payload.token);
        }
        if (data.payload?.refresh_token) {
            setRefreshToken(data.payload.refresh_token);
        }
        if (data && data.payload?.user_id) {
            localStorage.setItem("user_id", data.payload.user_id);
        }
        if (data && data.payload?.username) {
            localStorage.setItem("username", data.payload.username);
        }
        if (data && data.payload?.name) {
            localStorage.setItem("name", data.payload.name);
        }
        if (data && data.payload?.is_verified) {
            localStorage.setItem(
                "is_verified",
                data.payload.is_verified ? "true" : "false"
            );
        }
        addToast(
            "success",
            "User cookie validity checked successfully!",
            "Session Re-Validated!"
        );
        return data.payload?.token ?? "";
    }
}

async function checkTokenValidity(): Promise<boolean> {
    const token = getCookie("token");
    if (token) {
        return true;
    }

    const refresh_token = getCookie("refresh_token");
    if (!refresh_token) {
        return false;
    }

    addToast("info", "Refreshing session from cookie...", "Session Expired!");

    const newAccessToken = await refreshToken(refresh_token);
    if (newAccessToken !== "") {
        return true;
    } else {
        return false;
    }
}

export async function getValidToken(): Promise<string> {
    let token = getCookie("token");
    if (token) {
        return token;
    }

    const refresh_token = getCookie("refresh_token");
    if (!refresh_token) {
        return "";
    }

    token = await refreshToken(refresh_token);
    return token;
}
