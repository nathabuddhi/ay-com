import { API_URL } from "../env_var";
import type { ApiResponse } from "../types/api";
import type {
    Community,
    GetCategoriesResponse,
    GetCommunitiesResponse,
    GetCommunityMembersResponse,
} from "../types/community";
import { getValidToken } from "./token-controller";
import { returnDefaultError } from "./util";

export async function getCommunityCategories(): Promise<
    ApiResponse<GetCategoriesResponse>
> {
    try {
        const response = await fetch(`${API_URL}/community/getcategories`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
        });
        const data: ApiResponse<GetCategoriesResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<GetCategoriesResponse>(error);
    }
}

export async function getAllCommunities(): Promise<
    ApiResponse<GetCommunitiesResponse>
> {
    try {
        const response = await fetch(`${API_URL}/community/getallcommunities`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
        });
        const data: ApiResponse<GetCommunitiesResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<GetCommunitiesResponse>(error);
    }
}

export async function getCommunityById(
    id: string
): Promise<ApiResponse<Community>> {
    try {
        const response = await fetch(`${API_URL}/community/get/${id}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
        });
        const data: ApiResponse<Community> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<Community>(error);
    }
}

export async function getUserCommunities(): Promise<
    ApiResponse<GetCommunitiesResponse>
> {
    try {
        const response = await fetch(
            `${API_URL}/community/getusercommunities`,
            {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
            }
        );
        const data: ApiResponse<GetCommunitiesResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<GetCommunitiesResponse>(error);
    }
}

export async function getJoinedCommunities(): Promise<
    ApiResponse<GetCommunitiesResponse>
> {
    try {
        const response = await fetch(
            `${API_URL}/community/getusercommunities`,
            {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
            }
        );
        const data: ApiResponse<GetCommunitiesResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<GetCommunitiesResponse>(error);
    }
}

export async function getPendingJoinCommunities(): Promise<
    ApiResponse<GetCommunitiesResponse>
> {
    try {
        const response = await fetch(
            `${API_URL}/community/getuserpendingcommunities`,
            {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
            }
        );
        const data: ApiResponse<GetCommunitiesResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<GetCommunitiesResponse>(error);
    }
}

export async function createNewCommunity(
    formdata: FormData
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/community/create`, {
            method: "POST",
            headers: {
                Authorization: await getValidToken(),
            },
            body: formdata,
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function getCommunityMembers(
    community_id: string
): Promise<ApiResponse<GetCommunityMembersResponse>> {
    try {
        const response = await fetch(`${API_URL}/community/getmembers`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
                body: JSON.stringify({ value: community_id }),
            },
        });
        const data: ApiResponse<GetCommunityMembersResponse> =
            await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<GetCommunityMembersResponse>(error);
    }
}

export async function sendJoinRequest(
    community_id: string
): Promise<ApiResponse<number>> {
    try {
        const response = await fetch(`${API_URL}/community/join`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ community_id }),
        });
        const data: ApiResponse<number> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<number>(error);
    }
}

export async function acceptJoinRequest(
    community_id: string,
    user_id: string
): Promise<ApiResponse<number>> {
    try {
        const response = await fetch(`${API_URL}/community/accept`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ community_id, user_id }),
        });
        const data: ApiResponse<number> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<number>(error);
    }
}

export async function declineJoinRequest(
    community_id: string,
    user_id: string
): Promise<ApiResponse<number>> {
    try {
        const response = await fetch(`${API_URL}/community/decline`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ community_id, user_id }),
        });
        const data: ApiResponse<number> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<number>(error);
    }
}

export async function promoteMember(
    community_id: string,
    user_id: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/community/promote`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ community_id, user_id }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function demoteModerator(
    community_id: string,
    user_id: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/community/demote`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ community_id, user_id }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function searchCommunities(
    name: string,
    categories: string[]
): Promise<ApiResponse<GetCommunitiesResponse>> {
    try {
        const response = await fetch(`${API_URL}/community/search`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                name,
                categories,
            }),
        });
        const data: ApiResponse<GetCommunitiesResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<GetCommunitiesResponse>(error);
    }
}

export async function getTopMembers(
    community_id: string
): Promise<ApiResponse<GetCommunityMembersResponse>> {
    try {
        const response = await fetch(`${API_URL}/community/gettopmembers`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ value: community_id }),
        });
        const data: ApiResponse<GetCommunityMembersResponse> =
            await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<GetCommunityMembersResponse>(error);
    }
}
