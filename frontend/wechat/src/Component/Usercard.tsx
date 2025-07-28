'use client'
import { useRouter } from "next/navigation";
interface userCardProps {
    userID: string;
    username: string;
    email: string;
    onChatClick: (userID: string) => void;
}
function UserCard({ username, email, userID, onChatClick }: userCardProps) {
    const router = useRouter();
    return (
        <div  className="flex w-full max-w-[60rem] items-center justify-between p-6 mb-6 rounded-xl border border-sky-100 bg-sky-50 shadow-sm">
            <div>
                <p className="text-xl font-semibold text-sky-800">{username}</p>
                <p className="text-sm text-stone-500">{email}</p>
            </div>
            <button
                onClick={() => onChatClick(userID)}
                className="px-5 py-2 rounded-md text-sm font-medium border border-sky-700 text-sky-800 hover:bg-sky-700 hover:text-white transition duration-200"
            >
                Message
            </button>
        </div>
    );
}
export default UserCard