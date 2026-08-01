import { createContext, useContext, useState, useEffect} from "react"
import type {ReactNode,} from "react"

import { createClient } from '@supabase/supabase-js'
import type {Session} from '@supabase/supabase-js'

type UserData = {
    id: string
    email?: string
}

type AuthContextType = {
    user: UserData | null,
    authLoading: boolean,
    session: Session | null,
    logout: () => void,
}

type AuthProviderProps = {
  children: ReactNode;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined)

const supabase = createClient(
  import.meta.env.VITE_SUPABASE_URL,
  import.meta.env.VITE_SUPABASE_PUBLISHABLE_KEY
)
export function AuthProvider({children}: AuthProviderProps) {
    const [user, setUser] = useState<UserData | null>(null);
    const [session, setSession] = useState<Session | null>(null)
    const [authLoading, setLoading] = useState<boolean>(true)

    function login(user: UserData) {
        setUser(user)
    }

    function logout() {
        setUser(null)
    }

    function syncUserFromClaims(session: Session | null) {
        supabase.auth.getClaims().then(({ data, error }) => {
        if (error || !data) {
            logout()
            setLoading(false)
            return
        }
        setSession(session)

        login({
            id: data.claims.sub,
            email: data.claims.email,
        })

        setLoading(false)
        })
    }

    useEffect(() => {
        const { data: { subscription }, } = supabase.auth.onAuthStateChange((event, session) => {
            if (event === 'SIGNED_IN' || event === 'TOKEN_REFRESHED' || event === 'INITIAL_SESSION') {
                syncUserFromClaims(session)
            } else if (event == "SIGNED_OUT") {
                logout()
                setSession(null)
                setLoading(false)
                return
            }
        })

        return () => subscription.unsubscribe()
    }, [])


    return (
        <AuthContext.Provider
            value={{
                user,
                authLoading,
                session,
                logout
            }}
        >
        {children}
        </AuthContext.Provider>
    )
}


export function useAuth() {
    const context = useContext(AuthContext)

    if (!context) {
        throw new Error("Auth Context not intialized!")
    }

    return context
}