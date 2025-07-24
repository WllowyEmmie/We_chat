'use client'
import { useState } from "react"
import axios from "axios"
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import UserCard from "@/Component/Usercard"
interface UserType {
    id: string;
    username: string;
    email: string;
}
//testing
function Homepage() {
    const route = "http://localhost:8080/api/users";
    const [users, setUsers] = useState<UserType[]>([]);
    const router = useRouter();
    const [error, setError] = useState<string | null>(null);
    const [roomID, setRoomID] = useState<string | null>(null)
    const [currentUserID, setCurrentUserID] = useState<string | null>(null);
    const [members, setMembers] = useState<UserType[]>([])
    const API_URL = "http://localhost:8080/api";

    useEffect(() => {
        const fetchUsers = async () => {
            try {
                const response = await axios.get(route, {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem("token")}`

                    },
                });
                const userID = localStorage.getItem("userID");
                setCurrentUserID(userID);

                const data = response.data;

                setUsers(data.users);


            } catch (error: any) {

                setError(error.response?.data?.error || "Failed to fetch users");
            }
        };

        fetchUsers();


    }, []);



    const handleCreateRoom = async (otherUserID: string) => {
        if (otherUserID == currentUserID) {
            alert("Cannot create a room with yourself")
            return
        }
        console.log(otherUserID);
        try {
            const response = await axios.post(`${API_URL}/room`, {
                user1_id: localStorage.getItem("userID"),
                user2_id: otherUserID
            }, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem("token")}`
                }
            });
            const room = response.data.room;
            const roomName = response.data.room_name
            setRoomID(room.id)
            router.push(`chat/${room.id}`);
        } catch (error: any) {
            alert("Failed to make room")
            console.log("error :", error)
        }
    }


    return (
        <div className="size-full bg-slate-100 min-w-screen min-h-screen flex justify-center">
            <div className="h-full bg-slate-100 px-4 pt-6 m-4 bg-grey-100 rounded-lg shadow-sm">
                {users.map(user => (
                    <UserCard
                        key={user.id}
                        userID={user.id}
                        username={user.username}
                        email={user.email}
                        onChatClick={handleCreateRoom}
                    />
                ))}
            </div>
        </div>
    );
}
export default Homepage