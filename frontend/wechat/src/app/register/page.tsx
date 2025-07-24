'use client'
import Form from "@/Component/Form";
import { use } from "react";

function Register() {
    const route = process.env.NEXT_PUBLIC_API_URL
        ? `${process.env.NEXT_PUBLIC_API_URL}/register`
        : "http://localhost:8080/register";
    console.log(route)
    return (
        <Form routes={route} destination="/login" />
    )
}
export default Register