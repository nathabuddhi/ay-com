export interface AdminAllUsersResponse {
    users: AdminUserProfile[];
}

export interface AdminUserProfile {
    user_id: string;
    username: string;
    name: string;
    is_verified: boolean;
    is_banned: boolean;
}

export interface GetVerificationRequestsResponse {
    requests: VerifyAccountRequest[];
}

export interface VerifyAccountRequest {
    id: string;
    user_id: string;
    reason_text: string;
    status: string;
    submitted_at: string;
    selfie_url: string;
}

export interface UserReportsResponse {
    reports: UserReport[];
}

export interface UserReport {
    reporter_id: string;
    reported_id: string;
    reason: string;
}

export interface CommunityRequestsResponse {
    requests: CommunityRequest[];
}

export interface CommunityRequest {
    community_id: string;
    user_id: string;
    community_name: string;
    category: string;
    description: string;
    submitted_at: string;
    community_image_url: string;
    community_banner_url: string;
}
