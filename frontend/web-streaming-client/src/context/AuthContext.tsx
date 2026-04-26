import { createContext, useContext, useState } from "react";


interface AuthContextType {
    token: string | null
    login: (token: string) => void
    logout: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode}) {
    const [token, setToken] = useState<string | null>(
        localStorage.getItem("token")
    )

    const login = (token: string) => {
        localStorage.setItem("token", token)
        setToken(token)
    }

    const logout = () => {
        localStorage.removeItem("token")
        setToken(null)
    }
    return (
        <AuthContext.Provider value={{token, login, logout}}>
            {children}
        </AuthContext.Provider>
    )
}

export function useAuth() {
    const context = useContext(AuthContext)
    if (!context) throw new Error("useAuth must be used within AuthProvider")
        return context
}