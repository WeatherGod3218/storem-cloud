import { createContext, useContext, useState, useEffect} from "react"
import type {ReactNode} from "react"

import { createClient } from '@supabase/supabase-js'
import type {Session, AuthOtpResponse, OAuthResponse, AuthError, AuthResponse} from '@supabase/supabase-js'

type UserData = {
    id: string
    email?: string
}

type AuthContextType = {
    user: UserData | null,
    authLoading: boolean,
    session: Session | null,
    verifyOtp: (params: any) => Promise<AuthResponse>
    signOut: (params?: any) => Promise<{error: AuthError | null}>
    signInOtp: (params: any) => Promise<AuthOtpResponse>
    signInOAuth: (params: any) => Promise<OAuthResponse>

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


    async function signInOtp(params: any) {
        return supabase.auth.signInWithOtp(params)
    }

    async function signInOAuth(params: any) {
        return supabase.auth.signInWithOAuth(params)
    }

    async function signOut(params?: any) {
        return supabase.auth.signOut(params)
    }

    async function verifyOtp(params: any) {
        return supabase.auth.verifyOtp(params)
    }

    function syncUserFromClaims(session: Session | null) {
        supabase.auth.getClaims().then(({ data, error }) => {
        if (error || !data) {
            setUser(null)
            setLoading(false)
            return
        }
        setSession(session)

        setUser({
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
                setUser(null)                
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
                verifyOtp,
                signInOtp,
                signInOAuth,
                signOut
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