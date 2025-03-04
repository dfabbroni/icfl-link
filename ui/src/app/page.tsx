'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export default function Home() {
  const [selectedType, setSelectedType] = useState<'node' | 'researcher' | null>(null)
  const router = useRouter()

  const handleTypeSelection = (type: 'node' | 'researcher') => {
    setSelectedType(type)
  }

  const handleAction = (action: 'register' | 'login') => {
    if (selectedType === 'node') {
      router.push(action === 'register' ? '/nodes' : '/nodes/login')
    } else if (selectedType === 'researcher') {
      router.push(action === 'register' ? '/user' : '/user/login')
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <Card className="w-[350px]">
        <CardHeader>
          <CardTitle>Welcome</CardTitle>
          <CardDescription>Choose your role to continue</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex space-x-4">
              <Button
                onClick={() => handleTypeSelection('node')}
                variant={selectedType === 'node' ? 'default' : 'outline'}
                className="flex-1"
              >
                Node
              </Button>
              <Button
                onClick={() => handleTypeSelection('researcher')}
                variant={selectedType === 'researcher' ? 'default' : 'outline'}
                className="flex-1"
              >
                Researcher
              </Button>
            </div>
            {selectedType && (
              <div className="space-y-2">
                <Button onClick={() => handleAction('register')} className="w-full">
                  Register
                </Button>
                <Button onClick={() => handleAction('login')} variant="outline" className="w-full">
                  Login
                </Button>
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
