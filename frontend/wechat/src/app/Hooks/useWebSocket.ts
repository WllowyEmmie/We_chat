import { useEffect, useRef, useState } from "react";
export interface ClientMessage {
    type: "join" | "message";
    room: string;
    user: string;
    content?: string;
}
export interface ServerMessage {
    type: "message" | "history";
    room: string;
    user: {
        id: string;
        username: string;
        email?: string;
    };
    content?: string;
    messages?: ServerMessage[];
}

export function useWebSocket(roomId: string, userId: string) {
    const [messages, setMessages] = useState<any[]>([]);
    const wsRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        if (!roomId || !userId) return;
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            console.log("WebSocket already connected");
            return;
        }
        const WS_URL: string = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080/ws";
        console.log(WS_URL)
        const ws = new WebSocket(WS_URL);
        console.log(ws)
        wsRef.current = ws;

        ws.onopen = () => {
            console.log(" Connected to WebSocket");
            // Join the room
            ws.send(JSON.stringify({ type: "join", room: roomId, user: userId }));
        };

        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);

            if (data.type === "history") {
                // Replace entire history
                setMessages(data.messages);
            } else if (data.type === "message") {
                // Append live message
                setMessages(prev => [...prev, data.message]);
            }
        };

        ws.onerror = (err) => console.error("WebSocket error", err);

        ws.onclose = () => console.log("WebSocket closed");

        return () => {
            ws.close();
        };
    }, [roomId, userId]);

    const sendMessage = (content: string) => {
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            const msg: ClientMessage = { type: "message", room: roomId, user: userId, content };
            wsRef.current.send(JSON.stringify(msg));
        } else {
            console.warn("WebSocket is not open. Cannot send message yet.");
        }
    };

    return { messages, sendMessage };
}