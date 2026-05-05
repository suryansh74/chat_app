import { createContext, useContext, useEffect, useRef, useCallback, useState, type ReactNode } from "react"
import { useAuth } from "@/contexts/AuthContext"
import { createWebSocket } from "@/lib/api"

type MessageHandler = (data: unknown) => void

interface WebSocketContextType {
  subscribe: (type: string, handler: MessageHandler) => void
  unsubscribe: (type: string, handler: MessageHandler) => void
  isConnected: boolean
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined)

export function WebSocketProvider({ children }: { children: ReactNode }) {
  const { user } = useAuth()
  const wsRef = useRef<WebSocket | null>(null)
  const handlersRef = useRef<Map<string, Set<MessageHandler>>>(new Map())
  const [isConnected, setIsConnected] = useState(false)
  const initializedRef = useRef(false)

  const handleMessage = useCallback((event: MessageEvent) => {
    try {
      const data = JSON.parse(event.data)
      const type = data.type
      if (!type) return

      const typeHandlers = handlersRef.current.get(type)
      if (typeHandlers) {
        typeHandlers.forEach((handler) => handler(data))
      }

      const wildcardHandlers = handlersRef.current.get("*")
      if (wildcardHandlers) {
        wildcardHandlers.forEach((handler) => handler(data))
      }
    } catch {
    }
  }, [])

  useEffect(() => {
    if (!user?.id || initializedRef.current) return

    initializedRef.current = true
    const ws = createWebSocket(user.id)
    wsRef.current = ws

    ws.onopen = () => setIsConnected(true)
    ws.onmessage = handleMessage
    ws.onclose = () => {
      setIsConnected(false)
      wsRef.current = null
    }
    ws.onerror = () => setIsConnected(false)

    return () => {
      initializedRef.current = false
      ws.close()
      wsRef.current = null
      setIsConnected(false)
    }
  }, [user?.id, handleMessage])

  const subscribe = useCallback((type: string, handler: MessageHandler) => {
    const set = handlersRef.current.get(type) || new Set()
    set.add(handler)
    handlersRef.current.set(type, set)
  }, [])

  const unsubscribe = useCallback((type: string, handler: MessageHandler) => {
    const set = handlersRef.current.get(type)
    if (set) {
      set.delete(handler)
      if (set.size === 0) {
        handlersRef.current.delete(type)
      }
    }
  }, [])

  return (
    <WebSocketContext.Provider value={{ subscribe, unsubscribe, isConnected }}>
      {children}
    </WebSocketContext.Provider>
  )
}

export function useWebSocket() {
  const context = useContext(WebSocketContext)
  if (!context) {
    throw new Error("useWebSocket must be used within a WebSocketProvider")
  }
  return context
}
