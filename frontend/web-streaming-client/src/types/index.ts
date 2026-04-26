export interface Filem {
    ID: number
    Title: string
    Description: string
    Genre: string
    Year: number
    PosterURL: string
    Rating: number
    VideoURL: string
}

export interface User {
    username: string
    email: string
    password: string
}

export interface LoginResponse{
    token: string
}

export interface ApiResponse<T> {
    data?: T
    message?: string
    error?: string
}