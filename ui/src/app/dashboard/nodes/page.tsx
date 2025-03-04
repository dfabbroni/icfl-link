'use client'

import { useState, useEffect } from 'react'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { nodeService, Node } from '@/services/nodeService'

export default function NodesPage() {
  const [nodes, setNodes] = useState<Node[]>([])
  const [filter, setFilter] = useState<'all' | 'true' | 'false'>('all')

  useEffect(() => {
    fetchNodes()
  }, [])

  const fetchNodes = async () => {
    try {
      const fetchedNodes = await nodeService.getAll()
      setNodes(fetchedNodes)
    } catch (error) {
      console.error('Failed to fetch nodes:', error)
    }
  }

  const filteredNodes = nodes.filter(node => {
    if (filter === 'all') return true
    return node.Approved === (filter === 'true')
  })

  const handleApprove = async (id: string) => {
    if (confirm('Are you sure you want to approve this node?')) {
      try {
        await nodeService.approve(id)
        setNodes(nodes.map(node => 
          node.ID === id ? { ...node, Approved: true } : node
        ))
      } catch (error) {
        console.error('Failed to approve node:', error)
      }
    }
  }

  const handleReject = async (id: string) => {
    if (confirm('Are you sure you want to reject this node?')) {
      try {
        await nodeService.reject(id)
        setNodes(nodes.filter(node => node.ID !== id))
      } catch (error) {
        console.error('Failed to reject node:', error)
      }
    }
  }

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Nodes</h1>
      <div>
        <Select onValueChange={(value: 'all' | 'true' | 'false') => setFilter(value)}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Filter by approval" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All</SelectItem>
            <SelectItem value="true">Approved</SelectItem>
            <SelectItem value="false">Not Approved</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Username</TableHead>
            <TableHead>Approved</TableHead>
            <TableHead>Last Seen</TableHead>
            <TableHead>Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {filteredNodes.map((node) => (
            <TableRow key={node.ID}>
              <TableCell>{node.ID}</TableCell>
              <TableCell>{node.Username}</TableCell>
              <TableCell>{node.Approved ? 'Yes' : 'No'}</TableCell>
              <TableCell>{new Date(node.LastSeen).toLocaleString()}</TableCell>
              <TableCell>
                {!node.Approved && (
                  <div className="space-x-2">
                    <Button onClick={() => handleApprove(node.ID)}>Approve</Button>
                    <Button variant="destructive" onClick={() => handleReject(node.ID)}>Reject</Button>
                  </div>
                )}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
