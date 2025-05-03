export interface ApiResponse<T> {
    success: boolean;
    message: string;
    payload: T;
}

export interface StringPayload {
    value: string;
}