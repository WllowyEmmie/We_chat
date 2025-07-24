'use client'
import Form from "@/Component/Form";

function Login() {
    return (
        <Form routes="http://localhost:8080/login" destination="/home" />
    );
}
export default Login;