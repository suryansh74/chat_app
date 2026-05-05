import { useState, useEffect } from "react"
import { Search, X, MessageCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { searchApi, type MessageListItem, type FriendListItem } from "@/lib/api"

interface SearchPanelProps {
  onClose: () => void
  onSelectFriend?: (friendId: string) => void
}

type SearchMode = "friends" | "global" | "local"

export function SearchPanel({ onClose, onSelectFriend }: SearchPanelProps) {
  const [mode, setMode] = useState<SearchMode>("global")
  const [query, setQuery] = useState("")
  const [results, setResults] = useState<(FriendListItem | MessageListItem)[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    const search = async () => {
      if (!query.trim()) {
        setResults([])
        return
      }

      setLoading(true)

      if (mode === "friends") {
        const { data } = await searchApi.searchGlobal(query)
        if (data) {
          setResults(data.messages || [])
        }
      } else if (mode === "global") {
        const { data } = await searchApi.searchGlobal(query)
        if (data) {
          setResults(data.messages || [])
        }
      }

      setLoading(false)
    }

    const debounce = setTimeout(search, 300)
    return () => clearTimeout(debounce)
  }, [query, mode])

  const handleSelectMessage = (msg: MessageListItem) => {
    const friendId = msg.is_me ? msg.to_user_id : msg.from_user_id
    if (onSelectFriend) {
      onSelectFriend(friendId)
      onClose()
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="flex h-[80vh] w-full max-w-md flex-col rounded-lg bg-background shadow-lg">
        <div className="flex items-center justify-between border-b p-3">
          <h2 className="text-lg font-semibold">Search</h2>
          <Button variant="ghost" size="icon" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>

        <div className="flex gap-1 border-b p-2">
          <Button
            size="sm"
            variant={mode === "global" ? "default" : "ghost"}
            onClick={() => setMode("global")}
          >
            Global
          </Button>
          <Button
            size="sm"
            variant={mode === "friends" ? "default" : "ghost"}
            onClick={() => setMode("friends")}
          >
            Friends
          </Button>
        </div>

        <div className="relative p-3">
          <Search className="absolute left-5 top-5 h-4 w-4 text-muted-foreground" />
          <Input
            className="pl-9"
            placeholder="Search messages..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>

        <div className="flex-1 overflow-y-auto">
          {loading ? (
            <div className="flex justify-center py-4">
              <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
            </div>
          ) : results.length === 0 ? (
            <div className="p-4 text-center text-muted-foreground">
              {query ? "No results found" : "Start typing to search"}
            </div>
          ) : (
            <ul className="divide-y">
              {"friend_id" in results[0] ? (
                (results as FriendListItem[]).map((friend) => (
                  <li key={friend.friend_id}>
                    <button
                      onClick={() => {
                        if (onSelectFriend) {
                          onSelectFriend(friend.friend_id)
                          onClose()
                        }
                      }}
                      className="w-full p-3 text-left hover:bg-accent"
                    >
                      <p className="font-medium">{friend.friend_name}</p>
                      <p className="text-sm text-muted-foreground truncate">
                        {friend.last_message || "No messages"}
                      </p>
                    </button>
                  </li>
                ))
              ) : (
                (results as MessageListItem[]).map((msg) => (
                  <li key={msg.message_id}>
                    <button
                      onClick={() => handleSelectMessage(msg)}
                      className="w-full p-3 text-left hover:bg-accent"
                    >
                      <div className="flex items-center gap-2">
                        <MessageCircle className="h-4 w-4" />
                        <span className="text-sm font-medium">
                          {msg.is_me ? "You" : "Friend"}
                        </span>
                      </div>
                      <p className="text-sm text-muted-foreground truncate">
                        {msg.content}
                      </p>
                    </button>
                  </li>
                ))
              )}
            </ul>
          )}
        </div>
      </div>
    </div>
  )
}