import { useState, useEffect, useRef } from "react"
import { Send } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { chatApi, type MessageListItem, type FriendListItem, createWebSocket } from "@/lib/api"

interface ChatWindowProps {
  friend: FriendListItem
  userId: string
}

export function ChatWindow({ friend, userId }: ChatWindowProps) {
  const [messages, setMessages] = useState<MessageListItem[]>([])
  const [newMessage, setNewMessage] = useState("")
  const [loading, setLoading] = useState(true)
  const [ws, setWs] = useState<WebSocket | null>(null)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const initializedRef = useRef(false)

  const loadMessages = async () => {
    setLoading(true)
    const { data } = await chatApi.getMessages(friend.friend_id, 50)
    if (data?.messages) {
      setMessages(data.messages.reverse())
    }
    setLoading(false)
  }

  const initWebSocket = () => {
    const socket = createWebSocket(userId)
    socket.onopen = () => {
      console.log("WebSocket connected")
      setWs(socket)
    }
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data)
      if (data.type === "new_message") {
        const msg = {
          ...data.message,
          is_me: data.message.from_user_id === userId
        }
        if (!msg.is_me) {
          setMessages((prev) => [...prev, msg])
        }
      }
    }
    socket.onerror = (error) => {
      console.error("WebSocket error:", error)
    }
    socket.onclose = () => {
      console.log("WebSocket closed")
    }
  }

  useEffect(() => {
    if (initializedRef.current) return
    initializedRef.current = true

    loadMessages()
    initWebSocket()

    return () => {
      if (ws) {
        ws.close()
      }
      initializedRef.current = false
    }
  }, [friend.friend_id])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages])

  const sendMessage = async () => {
    if (!newMessage.trim()) return

    const { data } = await chatApi.sendMessage(friend.friend_id, newMessage)
    if (data?.message) {
      setMessages((prev) => [...prev, data.message])
      setNewMessage("")
    }
  }

  const formatTime = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
  }

  return (
    <div className="flex h-full flex-col">
      <div className="border-b p-3">
        <h2 className="font-semibold">{friend.friend_name}</h2>
        <p className="text-sm text-muted-foreground">{friend.friend_email}</p>
      </div>

      <div className="flex-1 overflow-y-auto p-3">
        {loading ? (
          <div className="flex items-center justify-center h-full">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>
        ) : messages.length === 0 ? (
          <div className="flex items-center justify-center h-full text-muted-foreground">
            No messages yet. Start the conversation!
          </div>
        ) : (
          <div className="space-y-2">
            {messages.map((msg) => (
              <div
                key={msg.message_id}
                className={`flex ${msg.is_me ? "justify-end" : "justify-start"}`}
              >
                <div
                  className={`max-w-[70%] rounded-lg p-2 ${
                    msg.is_me
                      ? "bg-primary text-primary-foreground"
                      : "bg-muted"
                  }`}
                >
                  <p>{msg.content}</p>
                  <p className="text-xs opacity-70">{formatTime(msg.created_at)}</p>
                </div>
              </div>
            ))}
            <div ref={messagesEndRef} />
          </div>
        )}
      </div>

      <div className="border-t p-3">
        <form
          onSubmit={(e) => {
            e.preventDefault()
            sendMessage()
          }}
          className="flex gap-2"
        >
          <Input
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Type a message..."
          />
          <Button type="submit" size="icon">
            <Send className="h-4 w-4" />
          </Button>
        </form>
      </div>
    </div>
  )
}