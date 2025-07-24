"use client";
import { useRouter, useParams } from "next/navigation";
import { useEffect, useState } from "react";
import axios from "axios";
import { useWebSocket, ServerMessage } from "@/app/Hooks/useWebSocket";
interface UserType {
  id: string;
  username: string;
  email: string;
}
function ChatRoomPage() {
  const router = useRouter();
  const params = useParams();
  const roomID = params.roomID as string;
  const [userID, setUserID] = useState<string | null>(null)
  const API_URL = "http://localhost:8080/api";
  const [members, setMembers] = useState<UserType[]>([])
  const [error, setError] = useState<string | null>(null);

  const currentUserID =
    typeof window !== "undefined" ? localStorage.getItem("userID") : null;

  const { messages, sendMessage } = useWebSocket(roomID, currentUserID || "");
  const [text, setText] = useState("");

  const handleSend = () => {
    if (text.trim()) {
      sendMessage(text);
      setText("");
    }
  };
  const handleKeyDOwn = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key == "Enter" && !e.shiftKey) {
      e.preventDefault();
      if (text.trim()) {
        sendMessage(text);
        setText("");
      }
    }
  }
  useEffect(() => {
    const fetchmembers = async (roomID: string) => {

      try {
        const response = await axios.get(`${API_URL}/room/${roomID}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`
          },
        });
        const data = response.data;
        setMembers(data.members || [])
        setUserID(localStorage.getItem("userID"))

      } catch (error: any) {
        setError(error.response?.data.error || "Failed to fetch room")
      }
    };
    fetchmembers(roomID);
  }, [])

  const otherUser = members?.find(m => m.id !== currentUserID);


  return (
    <div className="flex flex-col items-center min-h-screen bg-gradient-to-b from-sky-100 to-white p-6">

      <div className="w-full max-w-2xl flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold text-sky-800">
          {otherUser ? otherUser.username.toUpperCase() : "Loading..."}
        </h1>
        <span className="text-sm text-gray-500">Room ID: {roomID}</span>
      </div>

      <div className="w-full max-w-2xl bg-white rounded-xl shadow-md border border-gray-200 p-4 mb-4 h-[450px] overflow-y-auto flex flex-col space-y-2">
        {messages.length === 0 ? (
          <p className="text-gray-500 text-center italic">
            No messages yet... Start the conversation!
          </p>
        ) : (
          messages.map((msg: ServerMessage, idx: number) => {
            const isMine = msg.user.id === currentUserID;

            return (
              <div
                key={idx}
                className={`flex ${isMine ? "justify-end" : "justify-start"}`}
              >
                <div
                  className={`max-w-[70%] p-3 rounded-lg shadow-sm ${isMine
                    ? "bg-sky-600 text-white rounded-br-none"
                    : "bg-gray-200 text-gray-900 rounded-bl-none"
                    }`}
                >
                  <p className="text-xs font-semibold opacity-70">
                    {typeof msg.user === "object" ? msg.user.username : msg.user}
                  </p>
                  <p>{msg.content || "test"}</p>
                </div>
              </div>
            );
          })
        )}
      </div>

      <div className="flex w-full max-w-2xl gap-2">

        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={handleKeyDOwn}
          placeholder="Type a message..."
          className="flex-1 px-4 py-2 text-black border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-400"
        />
        <button
          onClick={handleSend}
          className="px-5 py-2 bg-sky-600 text-white rounded-lg shadow hover:bg-sky-700 transition"
        >
          Send
        </button>
      </div>
    </div>
  );
}
export default ChatRoomPage