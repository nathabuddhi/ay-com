export interface Community {
    community_id: string;
    community_name: string;
    creator_id: string;
    categories: string[];
    description: string;
    rules: string;
    created_at: string;
    icon_image: string;
    banner_image: string;
    member_count: number;
    role: string;
    is_pending: boolean;
}

export interface GetCommunitiesResponse {
    communities: Community[];
}

export interface GetCategoriesResponse {
    categories: string[];
}

export interface CommunityMember {
    user_id: string;
    role: string;
}

export interface GetCommunityMembersResponse {
    members: CommunityMember[];
}
