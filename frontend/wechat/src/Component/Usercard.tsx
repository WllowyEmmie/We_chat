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
        <div className=" flex w-lvh items-center justify-between p-4 mb-6 rounded-sm border border-sky-50 bg-sky-50 shadow-sm">
            <div className="space-y-0.5">
                <p className="text-lg font-semibold text-sky-800">{username}</p>
                <p className="font medium text-stone-500 text-sm">{email}</p>
            </div>
            <button onClick={() => onChatClick(userID)} className="w-1/3 py-2 my-1 border border-sky-700 text-sky-800 hover:border-transparent hover:bg-sky-700 hover:text-white active:bg-sky-700 rounded-sm">
                Message
            </button>
        </div>
    );
}
export default UserCard