import { useState, useEffect } from "react"
import { Search, UserPlus, X } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { searchApi, friendsApi, type Friend } from "@/lib/api"

interface FriendRequestModalProps {
  onClose: () => void
}

export function FriendRequestModal({ onClose }: FriendRequestModalProps) {
  const [query, setQuery] = useState("")
  const [results, setResults] = useState<Friend[]>([])
  const [selectedUser, setSelectedUser] = useState<Friend | null>(null)
  const [loading, setLoading] = useState(false)
  const [sending, setSending] = useState(false)
  const [sent, setSent] = useState(false)

  useEffect(() => {
    const searchUsers = async () => {
      if (query.length < 2) {
        setResults([])
        return
      }

      setLoading(true)
      const { data } = await searchApi.searchByEmail(query)
      if (data?.users) {
        setResults(data.users)
      }
      setLoading(false)
    }

    const debounce = setTimeout(searchUsers, 300)
    return () => clearTimeout(debounce)
  }, [query])

  const handleSendRequest = async () => {
    if (!selectedUser) return

    setSending(true)
    const { error } = await friendsApi.sendFriendRequest(selectedUser.friend_id)
    setSending(false)

    if (!error) {
      setSent(true)
      setTimeout(onClose, 1500)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="w-full max-w-md rounded-lg bg-background p-6 shadow-lg">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold">Add Friend</h2>
          <Button variant="ghost" size="icon" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>

        {sent ? (
          <div className="py-8 text-center">
            <UserPlus className="mx-auto h-12 w-12 text-green-500" />
            <p className="mt-4 font-medium">Friend request sent!</p>
          </div>
        ) : (
          <>
            <div className="relative mb-4">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                className="pl-9"
                placeholder="Search by name or email..."
                value={query}
                onChange={(e) => setQuery(e.target.value)}
              />
            </div>

            {loading ? (
              <div className="flex justify-center py-4">
                <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
              </div>
            ) : results.length > 0 ? (
              <ul className="mb-4 divide-y">
                {results.map((user) => (
                  <li key={user.id}>
                    <button
                      onClick={() => setSelectedUser(user)}
                      className={`w-full p-3 text-left hover:bg-accent ${
                        selectedUser?.id === user.id ? "bg-accent" : ""
                      }`}
                    >
                      <p className="font-medium">{user.friend_name}</p>
                      <p className="text-sm text-muted-foreground">{user.friend_email}</p>
                    </button>
                  </li>
                ))}
              </ul>
            ) : query.length >= 2 ? (
              <div className="py-4 text-center text-muted-foreground">
                No users found
              </div>
            ) : null}

            <Button
              className="w-full"
              disabled={!selectedUser || sending}
              onClick={handleSendRequest}
            >
              {sending ? (
                <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
              ) : (
                <>
                  <UserPlus className="mr-2 h-4 w-4" />
                  Send Friend Request
                </>
              )}
            </Button>
          </>
        )}
      </div>
    </div>
  )
}