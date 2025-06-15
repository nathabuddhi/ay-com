export interface AdminAllUsersResponse {
    users: AdminUserProfile[];
}

export interface AdminUserProfile {
    id: string;
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
    report_id: string;
    reporter_id: string;
    reported_id: string;
    reason: string;
    submitted_at: string;
    status: string;
}

export interface CommunityRequestsResponse {
    communities: CommunityRequest[];
}

export interface CommunityRequest {
    community_id: string;
    user_id: string;
    community_name: string;
    categories: string[];
    description: string;
    created_at: string;
    banner_image: string;
    icon_image: string;
}
