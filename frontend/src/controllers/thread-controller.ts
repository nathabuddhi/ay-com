import { API_URL } from "../env_var";
import type { ApiResponse } from "../types/api";
import type { Thread, ThreadResponse } from "../types/thread";
import type { UserProfile } from "../types/user";
import { getValidToken } from "./token-controller";
import { returnDefaultError } from "./user-controller";

export async function getThreadOwner(
    user_id: string
): Promise<UserProfile | null> {
    const response = await fetch(`${API_URL}/user/getuserid`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: await getValidToken(),
        },
        body: JSON.stringify({
            user_id,
        }),
    });

    const data: ApiResponse<UserProfile> = await response.json();
    if (!response.ok || !data) {
        return null;
    }

    return data.payload;
}

export async function getForYouThreads(): Promise<ApiResponse<ThreadResponse>> {
    try {
        const response = await fetch(`${API_URL}/thread/getallthreads`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({}),
        });
        const data: ApiResponse<ThreadResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<ThreadResponse>(error);
    }
}

export async function getFollowingThreads(): Promise<
    ApiResponse<ThreadResponse>
> {
    try {
        const response = await fetch(`${API_URL}/thread/getfollowingthreads`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({}),
        });
        const data: ApiResponse<ThreadResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<ThreadResponse>(error);
    }
}

export async function getThreadById(
    thread_id: string
): Promise<ApiResponse<Thread>> {
    try {
        const response = await fetch(`${API_URL}/thread/get/${thread_id}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
        });
        const data: ApiResponse<Thread> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<Thread>(error);
    }
}

export async function toggleLike(
    thread_id: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/thread/togglelike`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                thread_id,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function toggleRepost(
    thread_id: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/thread/togglerepost`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                thread_id,
                text: "",
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function toggleBookmark(
    thread_id: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/thread/togglebookmark`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                thread_id,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}