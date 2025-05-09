export interface VerifyAccountRequest {
    user_id: string;
    id: string;
    reason_text: string;
    status: string;
    submitted_at: string;
}

export interface UserSettings {
    font_size: string;
    font_color: string;
    is_private: boolean;
}

export interface UserProfile {
    user_id: string;
    username: string;
    name: string;
    bio: string;
    gender: string;
    is_verified: boolean;
    followers: number;
    following: number;
    date_of_birth: string;
    email: string;
    join_date: string;
}

export interface LoginResponse {
    token: string;
    user_id: string;
    username: string;
    name: string;
    is_verified: boolean;
}

export interface Settings {
    font_size: string;
    font_color: string;
    private: boolean;
}

export interface BlockedUser {
    user_id: string;
    username: string;
    name: string;
}
