'use client'
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import Form from "@/Component/Form";
export default function Home() {

  const router = useRouter()
  useEffect(() => {
    router.push("/register")
  }, [router])
  return null
}
