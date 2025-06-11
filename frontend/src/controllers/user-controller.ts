import type {
    ApiResponse,
    BlockedUserResponse,
    StringPayload,
} from "../types/api";
import type {
    GetFollowRecommendationsResponse,
    LoginResponse,
    Settings,
    UserProfile,
} from "../types/user";
import { API_URL } from "../env_var";
import { getValidToken, setRefreshToken, setToken } from "./token-controller";
import { returnDefaultError } from "./util";

export async function getUserById(
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

export async function tEMPLATE(user_id: string): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/${user_id}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({}),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function login(
    email: string,
    password: string
): Promise<ApiResponse<LoginResponse>> {
    try {
        const response = await fetch(`${API_URL}/user/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email: email,
                password: password,
            }),
        });
        const data: ApiResponse<LoginResponse> = await response.json();

        if (!response.ok || !data.success) {
            throw new Error(data.message || "Invalid email or password");
        }

        if (data && data.payload?.token) setToken(data.payload.token);
        if (data && data.payload?.refresh_token)
            setRefreshToken(data.payload.refresh_token);
        if (data && data.payload?.user_id) {
            localStorage.setItem("user_id", data.payload.user_id);
        }
        if (data && data.payload?.username) {
            localStorage.setItem("username", data.payload.username);
        }
        if (data && data.payload?.name) {
            localStorage.setItem("name", data.payload.name);
        }
        if (data && data.payload?.is_verified) {
            localStorage.setItem(
                "is_verified",
                data.payload.is_verified ? "true" : "false"
            );
        }
        if (data && data.payload?.is_admin) {
            localStorage.setItem(
                "is_admin",
                data.payload.is_admin ? "true" : "false"
            );
        }
        return data;
    } catch (error) {
        return returnDefaultError<LoginResponse>(error);
    }
}

export async function register(
    email: string,
    name: string,
    username: string,
    password: string,
    gender: string,
    date_of_birth: Date,
    security_question: string,
    security_answer: string,
    avatar: File,
    banner: File
): Promise<ApiResponse<null>> {
    try {
        const formData = new FormData();
        formData.append("email", email);
        formData.append("name", name);
        formData.append("username", username);
        formData.append("password", password);
        formData.append("gender", gender);
        formData.append(
            "date_of_birth",
            date_of_birth.toISOString().split("T")[0]
        );
        formData.append("security_question", security_question);
        formData.append("security_answer", security_answer);
        formData.append("avatar", avatar);
        formData.append("banner", banner);

        const response = await fetch(`${API_URL}/user/register`, {
            method: "POST",
            body: formData,
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function validateVerificationCode(
    email: string,
    code: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(
            `${API_URL}/user/validateverificationcode`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    email: email,
                    code: code,
                }),
            }
        );
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function requestVerificationCode(
    email: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(
            `${API_URL}/user/requestverificationcode`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    email: email,
                }),
            }
        );
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function getProfile(
    username: string
): Promise<ApiResponse<UserProfile>> {
    try {
        const response = await fetch(`${API_URL}/user/getprofile/` + username, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
        });
        const data: ApiResponse<UserProfile> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<UserProfile>(error);
    }
}

export async function getSelfProfile(): Promise<ApiResponse<UserProfile>> {
    try {
        const response = await fetch(`${API_URL}/user/getselfprofile`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({}),
        });
        const data: ApiResponse<UserProfile> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<UserProfile>(error);
    }
}

export async function updateProfile(
    name: string,
    username: string,
    bio: string,
    date_of_birth: Date,
    gender: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/user/updateprofile`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                name: name,
                username: username,
                bio: bio,
                date_of_birth: date_of_birth.toISOString().split("T")[0],
                gender: gender,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function getSecurityQuestion(
    email: string
): Promise<ApiResponse<StringPayload>> {
    try {
        const response = await fetch(`${API_URL}/user/getsecurityquestion`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email: email,
            }),
        });
        const data: ApiResponse<StringPayload> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<StringPayload>(error);
    }
}

export async function validateSecurityQuestion(
    email: string,
    answer: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/user/validatesecurityanswer`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email: email,
                answer: answer,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function changePassword(
    email: string,
    old_password: string,
    new_password: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/user/changepassword`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                email: email,
                old_password: old_password,
                new_password: new_password,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function resetPassword(
    email: string,
    code: string,
    new_password: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/user/resetpassword`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email: email,
                code: code,
                new_password: new_password,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function deactivateAccount(
    password: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/user/deactivateaccount`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                password: password,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function getSettings(): Promise<ApiResponse<Settings>> {
    try {
        const response = await fetch(`${API_URL}/user/getsettings`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({}),
        });
        const data: ApiResponse<Settings> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<Settings>(error);
    }
}

export async function updateSettings(
    font_size: string,
    font_color: string,
    is_private: boolean
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/user/updatesettings`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                font_size: font_size,
                font_color: font_color,
                private: is_private,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function getBlockedUsers(): Promise<
    ApiResponse<BlockedUserResponse>
> {
    try {
        const response = await fetch(`${API_URL}/user/getallblocked`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({}),
        });
        const data: ApiResponse<BlockedUserResponse> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<BlockedUserResponse>(error);
    }
}

export async function unblockUser(user_id: string): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/user/unblockuser`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                to_unblock_id: user_id,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function blockUser(user_id: string): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/user/blockuser`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                to_block_id: user_id,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function followUser(user_id: string): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/user/followuser`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                to_follow_id: user_id,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function unfollowUser(
    user_id: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/user/unfollowuser`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                to_unfollow_id: user_id,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function getFollowRecommendations(): Promise<
    ApiResponse<GetFollowRecommendationsResponse>
> {
    try {
        const response = await fetch(
            `${API_URL}/user/getfollowrecommendations`,
            {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
            }
        );
        const data: ApiResponse<GetFollowRecommendationsResponse> =
            await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<GetFollowRecommendationsResponse>(error);
    }
}

export async function reportUser(
    user_id: string,
    reason: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/user/report`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                reported_id: user_id,
                reason: reason,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}
