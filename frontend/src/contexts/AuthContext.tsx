import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  type ReactNode,
} from "react"
import { authApi } from "@/lib/api"

interface User {
  id: string
  name: string
  email: string
}

interface AuthContextType {
  user: User | null
  isVerified: boolean
  isLoading: boolean
  checkAuth: () => Promise<void>
  login: (email: string, password: string, rememberMe?: boolean) => Promise<{ error?: string }>
  register: (
    name: string,
    email: string,
    password: string,
    passwordConfirmation: string
  ) => Promise<{ error?: string }>
  logout: () => Promise<void>
  setUser: (user: User | null) => void
  setIsVerified: (verified: boolean) => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isVerified, setIsVerified] = useState(false)
  const [isLoading, setIsLoading] = useState(true)

  const checkAuth = useCallback(async () => {
    console.log("[AuthContext] checkAuth called")
    setIsLoading(true)
    
    console.log("[AuthContext] Calling /email_verification/verified...")
    const { data, error } = await authApi.checkVerified()
    console.log("[AuthContext] /email_verification/verified response:", { data, error })

    if (error || !data) {
      console.log("[AuthContext] checkVerified failed or no data, setting user=null, verified=false")
      setUser(null)
      setIsVerified(false)
      setIsLoading(false)
      return
    }

    console.log("[AuthContext] Calling /profile...")
    const { data: profileData, error: profileError } = await authApi.getProfile()
    console.log("[AuthContext] /profile response:", { profileData, profileError })

    if (profileError || !profileData) {
      console.log("[AuthContext] profile failed, setting user=null, verified=false")
      setUser(null)
      setIsVerified(false)
      setIsLoading(false)
      return
    }

    console.log("[AuthContext] Setting user and verified:", { user: profileData.user, verified: data.verified })
    setUser(profileData.user)
    setIsVerified(data.verified)
    setIsLoading(false)
  }, [])

   
  useEffect(() => {
    checkAuth()
  }, [checkAuth])

  const login = async (
    email: string,
    password: string,
    rememberMe?: boolean
  ): Promise<{ error?: string }> => {
    const { error } = await authApi.login({ email, password, remember_me: rememberMe })
    if (error) return { error }

    await checkAuth()
    return {}
  }

  const register = async (
    name: string,
    email: string,
    password: string,
    passwordConfirmation: string
  ): Promise<{ error?: string }> => {
    const { error } = await authApi.register({
      name,
      email,
      password,
      password_confirmation: passwordConfirmation,
    })
    if (error) return { error }

    await checkAuth()
    return {}
  }

  const logout = async () => {
    await authApi.logout()
    setUser(null)
    setIsVerified(false)
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isVerified,
        isLoading,
        checkAuth,
        login,
        register,
        logout,
        setUser,
        setIsVerified,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}