import { API_URL } from "../env_var";
import type {
    AdminAllUsersResponse,
    CommunityRequestsResponse,
    GetVerificationRequestsResponse,
    UserReportsResponse,
} from "../types/admin";
import type { ApiResponse } from "../types/api";

import { getValidToken } from "./token-controller";
import { returnDefaultError } from "./util";

export async function getAllUserVerificationRequests(): Promise<
    ApiResponse<GetVerificationRequestsResponse>
> {
    try {
        const response = await fetch(
            `${API_URL}/admin/getalluserverificationrequests`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
                body: JSON.stringify({}),
            }
        );
        const data: ApiResponse<GetVerificationRequestsResponse> =
            await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<GetVerificationRequestsResponse>(error);
    }
}

export async function approveUserVerificationRequest(
    id: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(
            `${API_URL}/admin/approveverificationrequest`,
            {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
                body: JSON.stringify({ value: id }),
            }
        );
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function rejectUserVerificationRequest(
    id: string,
    reason: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(
            `${API_URL}/admin/rejectverificationrequest`,
            {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
                body: JSON.stringify({ id: id, reason: reason }),
            }
        );
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function addThreadCategory(
    category: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/admin/addthreadcategory`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ value: category }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function deleteThreadCategory(
    category: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/admin/deletethreadcategory`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ value: category }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function getAllUsers(): Promise<
    ApiResponse<AdminAllUsersResponse>
> {
    try {
        const response = await fetch(`${API_URL}/admin/getallusers`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
        });
        const data: ApiResponse<AdminAllUsersResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<AdminAllUsersResponse>(error);
    }
}

export async function toggleUserBan(
    user_id: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/admin/toggleuserban`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ value: user_id }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}
export async function sendNewsletter(
    title: string,
    content: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/admin/sendnewsletter`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ title, content }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function getCommunityRequests(): Promise<
    ApiResponse<CommunityRequestsResponse>
> {
    try {
        const response = await fetch(`${API_URL}/admin/getcommunityrequests`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
        });
        const data: ApiResponse<CommunityRequestsResponse> =
            await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<CommunityRequestsResponse>(error);
    }
}

export async function approveCommunityRequest(
    requestId: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(
            `${API_URL}/admin/approvecommunityrequest`,
            {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
                body: JSON.stringify({ value: requestId }),
            }
        );
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function rejectCommunityRequest(
    requestId: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(
            `${API_URL}/admin/rejectcommunityrequest`,
            {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
                body: JSON.stringify({ value: requestId }),
            }
        );
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function getUserReports(): Promise<
    ApiResponse<UserReportsResponse>
> {
    try {
        const response = await fetch(`${API_URL}/admin/getuserreports`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
        });
        const data: ApiResponse<UserReportsResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<UserReportsResponse>(error);
    }
}

export async function approveReport(
    reportId: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/admin/approvereport`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ reportId }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function rejectReport(
    reportId: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/admin/rejectreport`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ reportId }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function getCommunityCategories(): Promise<ApiResponse<string[]>> {
    try {
        const response = await fetch(
            `${API_URL}/admin/getcommunitycategories`,
            {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
            }
        );
        const data: ApiResponse<string[]> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<string[]>(error);
    }
}

export async function addCommunityCategory(
    category: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/admin/addcommunitycategory`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({ category }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function deleteCommunityCategory(
    category: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(
            `${API_URL}/admin/deletecommunitycategory`,
            {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
                body: JSON.stringify({ category }),
            }
        );
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}
