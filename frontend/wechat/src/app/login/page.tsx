'use client'
import Form from "@/Component/Form";

function Login() {
    const route = process.env.NEXT_PUBLIC_API_URL
        ? `${process.env.NEXT_PUBLIC_API_URL}/login`
        : "http://localhost:8080/login";
    console.log(route)
    return (

        <Form routes={route} destination="/home" />
    );
}
export default Login;