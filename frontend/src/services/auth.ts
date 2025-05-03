import type { User } from "../types/user";

export function isLoggedIn(): boolean {
    return !!getToken();
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

export function getCurrentUser(): User | null {
    const token = getToken();
    if (!token) return null;

    try {
        return JSON.parse(localStorage.getItem("user") || "null");
    } catch (error) {
        console.error("Error parsing user data", error);
        return null;
    }
}

export function setCurrentUser(user: User): void {
    localStorage.setItem("user", JSON.stringify(user));
}

export function logout(): void {
    removeToken();
    localStorage.removeItem("user");
    window.location.href = "/";
}
