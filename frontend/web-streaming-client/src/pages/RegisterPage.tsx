import { useState } from "react"
import { useNavigate } from "react-router-dom"
import api from "../api"

export const RegisterPage = () => {
    const [username, setUsername] = useState("")
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")    
    const [error, setError] = useState("")
    
    // 1. Tambahkan state untuk show/hide password
    const [showPassword, setShowPassword] = useState(false)
    
    const navigate = useNavigate()

 const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
        await api.post("/register", { username, email, password })
        navigate("/login")
    } catch (err) {
        setError("Register gagal, coba lagi")
    }
}

    return (
        <div>
            <form onSubmit={handleSubmit}>
                <p>Username</p>
                <input type="username" value={username} onChange={e => setUsername(e.target.value)} />
                <br />
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

export default RegisterPage
