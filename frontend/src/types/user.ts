export interface AuthResponse {
    success: boolean;
    message: string;
    payload: {
        token?: string;
        user?: User;
    };
}

export interface User {
    id: string;
    username: string;
    email: string;
}

export interface UserCredentials {
    email: string;
    password: string;
}
