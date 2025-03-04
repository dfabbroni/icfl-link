'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { nodeService } from '@/services/nodeService'
import { ApiError } from '@/services/api'

export default function NodeLogin() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const router = useRouter()

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
      const response = await nodeService.login({ Username: username, Password: password })
      const status = response.Approved ? 'approved' : 'pending'
      router.push(`/nodes/status?approved=${status}&username=${encodeURIComponent(username)}`)
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
          <CardTitle>Node Login</CardTitle>
          <CardDescription>Check your node approval status</CardDescription>
        </CardHeader>
        <CardContent>
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
              {isLoading ? 'Checking status...' : 'Check Status'}
            </Button>
          </form>
        </CardContent>
        <CardFooter className="flex justify-between">
          <Link href="/nodes" className="text-sm text-blue-500 hover:underline">
            Register as Node
          </Link>
          <Link href="/user/login" className="text-sm text-blue-500 hover:underline">
            Login as Researcher
          </Link>
        </CardFooter>
      </Card>
    </div>
  )
}
