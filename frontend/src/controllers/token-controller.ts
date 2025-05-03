import type { BoolPayload } from "../types/api";

export async function isLoggedIn(): Promise<boolean> {
    if (!!getToken()) return false;

    return await checkTokenValidity();
}

export function getToken(): string | null {
    return localStorage.getItem("jwt");
}

export function setToken(token: string): void {
    localStorage.setItem("jwt", token);
}

export function removeToken(): void {
    localStorage.removeItem("jwt");
}

export function logout(): void {
    removeToken();
    window.location.href = "/";
}

async function checkTokenValidity(): Promise<boolean> {
    const token = getToken();
    if (!token) return false;
    const response = await fetch(
        "http://localhost:5000/user/checktoken/" + token,
        {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({}),
        }
    );
    const data: BoolPayload = await response.json();

    if (!response.ok || !data.value) {
        removeToken();
        return false;
    }

    return true;
}
