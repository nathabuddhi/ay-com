import { API_URL } from "../env_var";
import type { ApiResponse } from "../types/api";
import { getValidToken } from "./token-controller";
import { returnDefaultError } from "./util";

export async function getCommunityCategories(): Promise<ApiResponse<string[]>> {
    try {
        const response = await fetch(`${API_URL}/community/getcategories`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
        });
        const data: ApiResponse<string[]> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<string[]>(error);
    }
}
