import type { UserProfile } from "../types/user";
import { getToken } from "./token-controller";

async function getThreadPoster(user_id: string): Promise<UserProfile> {
    const token = getToken();

    const response = await fetch("http://localhost:5000/user/getthreadowner", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: token,
        },
        body: JSON.stringify({
            user_id,
        }),
    });

    const data: UserProfile = await response.json();
    if (!response.ok || !data) {
        return null;
    }

    return true;
}