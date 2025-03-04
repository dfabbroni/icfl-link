'use client'

import { useAuth } from '@/hooks/useAuth'

export default function Dashboard() {
  const { isLoading } = useAuth()

  if (isLoading) {
    return <div>Loading...</div>
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Dashboard</h1>
      <p>Welcome to your dashboard. Select an option from the sidebar to get started.</p>
    </div>
  )
}
