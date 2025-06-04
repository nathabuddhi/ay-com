import { API_URL } from "../env_var";
import type { ApiResponse } from "../types/api";
import type { ThreadDetailResponse, ThreadResponse } from "../types/thread";
import { getValidToken } from "./token-controller";
import { returnDefaultError } from "./util";

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
): Promise<ApiResponse<ThreadDetailResponse>> {
    try {
        const response = await fetch(`${API_URL}/thread/get/${thread_id}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
        });
        const data: ApiResponse<ThreadDetailResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<ThreadDetailResponse>(error);
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

export async function getBookmarkedThreads(): Promise<
    ApiResponse<ThreadResponse>
> {
    try {
        const response = await fetch(`${API_URL}/thread/getbookmarks`, {
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

export async function getRepostedThreads(
    user_id: string
): Promise<ApiResponse<ThreadResponse>> {
    try {
        const response = await fetch(`${API_URL}/thread/getreposts`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ user_id }),
        });
        const data: ApiResponse<ThreadResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<ThreadResponse>(error);
    }
}

export async function postThread(
    content: string,
    category: string,
    reply_permission: string,
    media: FileList | null,
    poll: string[] | null,
    reply: string
): Promise<ApiResponse<ThreadResponse>> {
    try {
        const formData = new FormData();
        formData.append("content", content);
        formData.append("category", category);
        formData.append("reply_permission", reply_permission);
        formData.append("reply_to", reply || "");
        formData.append("media_count", media?.length.toString() ?? "0");
        if (media) {
            for (let i = 0; i < media.length; i++) {
                formData.append(`media_${i}`, media[i]);
            }
        }

        formData.append("poll_count", poll?.length.toString() ?? "0");
        poll?.forEach((option, index) => {
            formData.append(`poll_${index}`, option);
        });

        const response = await fetch(`${API_URL}/thread/create`, {
            method: "POST",
            headers: {
                Authorization: await getValidToken(),
            },
            body: formData,
        });
        const data: ApiResponse<ThreadResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<ThreadResponse>(error);
    }
}

export async function getUserThreads(
    user_id: string
): Promise<ApiResponse<ThreadResponse>> {
    try {
        const response = await fetch(`${API_URL}/thread/getuserthreads`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ user_id }),
        });
        const data: ApiResponse<ThreadResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<ThreadResponse>(error);
    }
}

export async function getUserLikedThreads(): Promise<
    ApiResponse<ThreadResponse>
> {
    try {
        const response = await fetch(`${API_URL}/thread/getuserlikedthreads`, {
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

export async function getUserMediaThreads(
    user_id: string
): Promise<ApiResponse<ThreadResponse>> {
    try {
        const response = await fetch(`${API_URL}/thread/getusermediathreads`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ user_id }),
        });
        const data: ApiResponse<ThreadResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<ThreadResponse>(error);
    }
}

export async function getUserReplies(
    user_id: string
): Promise<ApiResponse<ThreadResponse>> {
    try {
        const response = await fetch(`${API_URL}/thread/getuserreplies`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ user_id }),
        });
        const data: ApiResponse<ThreadResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<ThreadResponse>(error);
    }
}

export async function voteThreadPoll(
    thread_id: string,
    option: string
): Promise<ApiResponse<ThreadResponse>> {
    try {
        const response = await fetch(`${API_URL}/thread/vote`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                thread_id,
                content: option,
            }),
        });
        const data: ApiResponse<ThreadResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<ThreadResponse>(error);
    }
}

export async function replyToThread(
    thread_id: string,
    comment: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/thread/reply`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                thread_id,
                content: comment,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}
