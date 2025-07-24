'use client'
import React, { useState } from "react";
import axios from "axios";
import { useRouter } from "next/navigation";
interface Form {
    username: string;
    email?: string;
    password: string
}
type FormProps = {
    routes: string,
    destination: string,
}
function Form({ routes, destination }: FormProps) {
    const [form, setForm] = useState<Form>({
        username: "",
        email: "",
        password: "",
    });
    const router = useRouter()
    const [error, setError] = useState<string | null>(null)
    const validateForm = () => {
        if (routes.endsWith("/register")) {
            if (!form.email?.trim() || !form.username.trim() || !form.password.trim()) {
                return ("All fields are required");
            }
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(form.email)) {
                return ("Input a valid email");
            }

        }
        const passwordRegex = /^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*]).{6,}$/;
        if (!passwordRegex.test(form.password)) {
            return ("Password Should have a Capital letter a number and a symbol and must be at leats 6 characters");
        }

    };
    const handleSubmit = async (event: React.FormEvent) => {
        event.preventDefault();
        const validateError = validateForm();
        if (validateError) {
            setError(validateError)
            return
        }
        try {
            const response = await axios.post(routes, form);
            if (routes.endsWith("/login")) {
                const data = response.data;
                const token = data.token;
                const user = data.user;
                localStorage.setItem("token", token);
                localStorage.setItem("userID", user.id);
            }
            router.push(destination as string)
        } catch (error: any) {
            setError(error.response?.data?.error || "Failed")
        }
    }
    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-300">

            <form onSubmit={handleSubmit} className="bg-sky-50 p-8 rounded-lg shadow-md w-full max-w-md">
                <h2 className="text-2xl font-bold mb-6 text-center text-sky-800">
                    {routes.endsWith("/register") ? "Register" : "Login"}
                </h2>
                {error && (
                    <div className="text-center bg-red-100 text-red-600 p-2 mb-4 rounded text-sm">
                        {error}
                    </div>
                )}
                {routes.endsWith("/register") ?
                    <div className="mb-4">
                        <label htmlFor="username" className="block mb-1 text-sm font-medium text-sky-800">Username</label>
                        <input type="text" id="username" value={form.username} onChange={(e) => setForm({ ...form, username: e.target.value })}
                            className="w-full px-3 text-black py-2 border rounded focus:outline-none focus:ring focus:border-blue-300" required />
                    </div> : null}
                <div className="mb-4">
                    <label htmlFor="email" className="block mb-1 text-sm font-medium text-sky-800">Email</label>
                    <input type="email" id="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })}
                        className="w-full px-3 py-2 text-black border rounded focus:outline-none focus:ring focus:border-blue-300" required />
                </div>

                <div className="mb-4">
                    <label htmlFor="password" className="block mb-1 text-sm font-medium text-sky-800">Password</label>
                    <input type="password" id="password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })}
                        className="w-full px-3 text-black py-2 border rounded focus:outline-none focus:ring focus:border-blue-300" />
                </div>
                <button className="w-full py-2 my-1 border border-sky-700 text-sky-700 hover:border-transparent hover:bg-sky-700 hover:text-white active:bg-sky-700 rounded-sm">
                    {routes.endsWith("/register") ? "Register" : "Login"}
                </button>
                {routes.endsWith("/register") ? <button onClick={() => router.push('/login')} className="w-full py-2 my-1 border border-sky-700 text-sky-700 hover:border-transparent hover:bg-sky-700 hover:text-white active:bg-sky-700 rounded-sm">
                    Already have an account? Login
                </button> : null}
            </form>
        </div>
    );
}
export default Form;