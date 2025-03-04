'use client'

import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

export default function NodeStatus() {
  const [isApproved, setIsApproved] = useState<boolean | null>(null)
  const [username, setUsername] = useState<string | null>(null)
  const router = useRouter()
  const searchParams = useSearchParams()

  useEffect(() => {
    const approved = searchParams.get('approved')
    const user = searchParams.get('username')
    
    if (approved === null || user === null) {
      router.push('/nodes/login')
    } else {
      setIsApproved(approved === 'approved')
      setUsername(user)
    }
  }, [searchParams, router])

  const handleLogout = () => {
    // Clear any stored credentials if necessary
    router.push('/nodes/login')
  }

  if (isApproved === null || username === null) {
    return <div>Loading...</div>
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <Card className="w-[350px]">
        <CardHeader>
          <CardTitle>Node Status</CardTitle>
          <CardDescription>Your current approval status</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="mb-4">Username: {username}</p>
          <p className={`text-lg font-semibold ${isApproved ? 'text-green-500' : 'text-yellow-500'}`}>
            Your node is currently {isApproved ? 'approved' : 'pending approval'}.
          </p>
        </CardContent>
        <CardFooter>
          <Button onClick={handleLogout} className="w-full">Logout</Button>
        </CardFooter>
      </Card>
    </div>
  )
}
