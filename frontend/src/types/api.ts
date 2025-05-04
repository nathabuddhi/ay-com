export interface ApiResponse<T> {
    success: boolean;
    message: string;
    payload: T | null;
}

export interface StringPayload {
    value: string;
}

export interface BoolPayload {
    value: boolean;
}

export interface NumberPayload {
    value: number;
}


