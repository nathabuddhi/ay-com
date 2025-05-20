export interface Notification {
    notification_id: string;
    title: string;
    content: string;
    from: string;
    timestamp: Date;
    read: boolean;
}

export interface GetNotificiationResponse {
    notifications: Notification[];
}
