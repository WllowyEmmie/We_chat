'use client'
import Form from "@/Component/Form";
import { use } from "react";

function Register() {
    return (
        <Form routes="http://localhost:8080/register" destination="/login" />
    )
}
export default Register