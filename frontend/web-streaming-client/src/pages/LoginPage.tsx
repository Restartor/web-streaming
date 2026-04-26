import { useState } from "react"
import { useAuth } from "../context/AuthContext"
import { useNavigate } from "react-router-dom"
import api from "../api"

export const LoginPage = () => {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")    
    const [error, setError] = useState("")
    
    // 1. Tambahkan state untuk show/hide password
    const [showPassword, setShowPassword] = useState(false)
    
    const { login } = useAuth()
    const navigate = useNavigate()

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        try {
            const res = await api.post("/login", {email, password})
            login(res.data.token)
            navigate("/")
        } catch (err) {
            setError("email atau password salah")
        }
    }

    return (
        <div>
            <form onSubmit={handleSubmit}>
                <p>Email</p>
                <input type="email" value={email} onChange={e => setEmail(e.target.value)} />
                <br />
                
                <p>Password</p>
                {/* 2. Ubah type berdasarkan state showPassword */}
                <input 
                    type={showPassword ? "text" : "password"} 
                    value={password} 
                    onChange={e => setPassword(e.target.value)}
                />
                
                {/* 3. Tambahkan checkbox atau button untuk toggle */}
                <div style={{ marginTop: "5px" }}>
                    <input 
                        type="checkbox" 
                        id="show" 
                        onChange={() => setShowPassword(!showPassword)} 
                    />
                    <label htmlFor="show"> Lihat Password</label>
                </div>

                {error && <p style={{ color: "red" }}>{error}</p>}
                <button type="submit">Login</button>
            </form>
        </div>
    )
}

export default LoginPage
