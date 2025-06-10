import { API_URL } from "../env_var";
import type { ApiResponse } from "../types/api";
import type { GetVerificationRequestsResponse } from "../types/user";
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
