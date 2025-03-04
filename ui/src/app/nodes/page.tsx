'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { nodeService } from '@/services/nodeService'
import { ApiError } from '@/services/api'

export default function NodeRegister() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [publicKey, setPublicKey] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const router = useRouter()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setIsLoading(true)

    if (!username || !password || !confirmPassword || !publicKey) {
      setError('Please fill in all fields')
      setIsLoading(false)
      return
    }

    if (password !== confirmPassword) {
      setError('Passwords do not match')
      setIsLoading(false)
      return
    }

    if (password.length < 8) {
      setError('Password must be at least 8 characters long')
      setIsLoading(false)
      return
    }

    try {
      console.log(username, password, publicKey)
      const response = await nodeService.register({ Username: username, Password: password, PublicKey: publicKey })
      router.push(`/nodes/status?approved=${response.Approved}&username=${encodeURIComponent(username)}`)
    } catch (error) {
      if (error instanceof ApiError) {
        setError(`Registration failed: ${error.message}`)
      } else {
        console.error('Registration failed:', error)
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
          <CardTitle>Node Registration</CardTitle>
          <CardDescription>Create a new node account</CardDescription>
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
            <Input
              type="password"
              placeholder="Confirm Password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
            />
            <Input
              type="text"
              placeholder="Public Key"
              value={publicKey}
              onChange={(e) => setPublicKey(e.target.value)}
              required
            />
            {error && <p className="text-red-500 text-sm">{error}</p>}
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? 'Registering...' : 'Register'}
            </Button>
          </form>
        </CardContent>
        <CardFooter className="flex justify-between">
          <Link href="/nodes/login" className="text-sm text-blue-500 hover:underline">
            Login as Node
          </Link>
          <Link href="/user" className="text-sm text-blue-500 hover:underline">
            Register as Researcher
          </Link>
        </CardFooter>
      </Card>
    </div>
  )
}
