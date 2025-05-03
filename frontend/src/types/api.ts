export interface ApiResponse<T> {
    Success: boolean;
    message: string;
    payload: T;
}
