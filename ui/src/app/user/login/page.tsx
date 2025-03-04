'use client'

import { useState, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { userService } from '@/services/userService'
import { ApiError } from '@/services/api'

export default function UserLogin() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const router = useRouter()
  const searchParams = useSearchParams()

  useEffect(() => {
    if (searchParams.get('registered') === 'true') {
      setMessage('Registration successful. Please wait for account approval before logging in.')
    }
  }, [searchParams])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setIsLoading(true)

    if (!username || !password) {
      setError('Please fill in all fields')
      setIsLoading(false)
      return
    }

    try {
      const { token } = await userService.login({ Username: username, Password: password })
      if (token) {
        localStorage.setItem('authToken', token)
        router.push('/dashboard')
      } else {
        setError('Your account is not yet approved. Please try again later.')
      }
    } catch (error) {
      if (error instanceof ApiError) {
        setError(`Login failed: ${error.message}`)
      } else {
        console.error('Login failed:', error)
        setError('An unexpected error occurred. Please try again.')
      }
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <Card className="w-[350px]">
        <CardHeader>
          <CardTitle>Researcher Login</CardTitle>
          <CardDescription>Access your researcher account</CardDescription>
        </CardHeader>
        <CardContent>
          {message && <p className="text-green-500 text-sm mb-4">{message}</p>}
          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              type="text"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
            <Input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            {error && <p className="text-red-500 text-sm">{error}</p>}
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? 'Logging in...' : 'Login'}
            </Button>
          </form>
        </CardContent>
        <CardFooter className="flex justify-between">
          <Link href="/user" className="text-sm text-blue-500 hover:underline">
            Register as Researcher
          </Link>
          <Link href="/nodes/login" className="text-sm text-blue-500 hover:underline">
            Login as Node
          </Link>
        </CardFooter>
      </Card>
    </div>
  )
}
