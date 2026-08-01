import { useState, useEffect } from 'react'

import { createClient } from '@supabase/supabase-js'
import { useAuth } from '@/context/AuthContext'

const supabase = createClient(
  import.meta.env.VITE_SUPABASE_URL,
  import.meta.env.VITE_SUPABASE_PUBLISHABLE_KEY
)

export default function LoginPage() {
  const {user, logout} = useAuth()

  const [loading, setLoading] = useState(false)
  const [email, setEmail] = useState('')

  const params = new URLSearchParams(window.location.search)
  const hasTokenHash = params.get('token_hash')

  const [verifying, setVerifying] = useState(!!hasTokenHash)
  const [authError, setAuthError] = useState<string | null>(null)
  const [authSuccess, setAuthSuccess] = useState(false)

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const token_hash = params.get('token_hash')
    const type = params.get('type')

    if (token_hash) {
      // Verify the OTP token
      supabase.auth
        .verifyOtp({
          token_hash,
          type: type || 'email',
        })
        .then(({ error }) => {
          if (error) {
            setAuthError(error.message)
          } else {
            setAuthSuccess(true)
            // Clear URL params
            window.history.replaceState({}, document.title, '/')
          }
          setVerifying(false)
        })
    }
  }, [])

  const handleLogin = async (event: any) => {
    event.preventDefault()
    setLoading(true)
    const { error } = await supabase.auth.signInWithOtp({
      email,
      options: {
        emailRedirectTo: window.location.origin,
      },
    })
    if (error) {
      alert(error.message)
    } else {
      alert('Check your email for the login link!')
    }
    setLoading(false)
  }

  const handleLoginGoogle = async (event: any) => {
    event.preventDefault()
    setLoading(true)
    const { error } = await supabase.auth.signInWithOAuth({
      provider: "google",
    })
    if (error) {
      alert(error.message)
    } else {
      alert('Signed in with google!')
    }
    setLoading(false)
  }

  const handleLogout = async () => {
    await supabase.auth.signOut()
    logout()
  }

  // Show verification state
  if (verifying) {
    return (
      <div>
        <h1>Authentication</h1>
        <p>Confirming your magic link...</p>
        <p>Loading...</p>
      </div>
    )
  }

  // Show auth error
  if (authError) {
    return (
      <div>
        <h1>Authentication</h1>
        <p>✗ Authentication failed</p>
        <p>{authError}</p>
        <button
          onClick={() => {
            setAuthError(null)
            window.history.replaceState({}, document.title, '/')
          }}
        >
          Return to login
        </button>
      </div>
    )
  }

  // Show auth success (briefly before claims load)
  if (authSuccess && !user) {
    return (
      <div>
        <h1>Authentication</h1>
        <p>✓ Authentication successful!</p>
        <p>Loading your account...</p>
      </div>
    )
  }

  // If user is logged in, show welcome screen
  if (user) {
    return (
      <div>
        <h1>Welcome!</h1>
        <p>You are logged in as: {user.email}</p>
        <button onClick={handleLogout}>Sign Out</button>
      </div>
    )
  }

  // Show login form
  return (
    <div>
      <h1>Supabase + React</h1>
      <p>Sign in via magic link with your email below</p>
      <form onSubmit={handleLogin}>
        <input
          type="email"
          placeholder="Your email"
          value={email}
          required={true}
          onChange={(e) => setEmail(e.target.value)}
        />
        <button disabled={loading}>
          {loading ? <span>Loading</span> : <span>Send magic link</span>}
        </button>
      </form>
      <br/>
      <p>Sign in via magic link with your google account</p>
      <form onSubmit={handleLoginGoogle}>
        <button disabled={loading}>
          {loading ? <span>Loading</span> : <span>Login with Google</span>}
        </button>
      </form>
    </div>
  )
}