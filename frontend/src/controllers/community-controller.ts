import { API_URL } from "../env_var";
import type { ApiResponse } from "../types/api";
import type { GetCategoriesResponse } from "../types/community";
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
