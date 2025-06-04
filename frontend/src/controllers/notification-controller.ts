import type { ApiResponse } from "../types/api";
import { API_URL } from "../env_var";
import { getValidToken } from "./token-controller";
import type { GetNotificiationResponse } from "../types/notification";
import { returnDefaultError } from "./util";
import type { NotificationSettings } from "../types/settings";

export async function getAllNotifications(): Promise<
    ApiResponse<GetNotificiationResponse>
> {
    try {
        const response = await fetch(
            `${API_URL}/notification/getallnotifications`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
                body: JSON.stringify({}),
            }
        );
        const data: ApiResponse<GetNotificiationResponse> =
            await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<GetNotificiationResponse>(error);
    }
}

export async function markNotifAsRead(
    notifId: string
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(
            `${API_URL}/notification/marknotificationasread`,
            {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
                body: JSON.stringify({ notification_id: notifId }),
            }
        );
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function deleteNotif(notifId: string): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(
            `${API_URL}/notification/deletenotification`,
            {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
                body: JSON.stringify({ notification_id: notifId }),
            }
        );
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function clearNotifs(): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(
            `${API_URL}/notification/clearnotifications`,
            {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: await getValidToken(),
                },
                body: JSON.stringify({}),
            }
        );
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}

export async function getNotificationSettings(): Promise<
    ApiResponse<NotificationSettings>
> {
    try {
        const response = await fetch(`${API_URL}/notification/getsettings`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({}),
        });
        const data: ApiResponse<NotificationSettings> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<NotificationSettings>(error);
    }
}

export async function updateNotificationSettings(
    notifSettings: NotificationSettings
): Promise<ApiResponse<null>> {
    try {
        const response = await fetch(`${API_URL}/notification/updatesettings`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                Authorization: await getValidToken(),
            },
            body: JSON.stringify({
                notif_like: notifSettings.notif_like,
                notif_repost: notifSettings.notif_repost,
                notif_follow: notifSettings.notif_follow,
                notif_mention: notifSettings.notif_mention,
                notif_community: notifSettings.notif_community,
                notif_newsletter: notifSettings.notif_newsletter,
            }),
        });
        const data: ApiResponse<null> = await response.json();

        return data;
    } catch (error) {
        return returnDefaultError<null>(error);
    }
}
