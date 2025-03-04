'use client'

import { useAuth } from '@/hooks/useAuth'
import { useEffect, useState } from 'react'
import { nodeService, Node } from '@/services/nodeService'
import { metadataService, Metadata } from '@/services/metadataService'
import { experimentService, Experiment } from '@/services/experimentService'
import { FaDatabase, FaFlask } from 'react-icons/fa'
import { FaCircleNodes } from "react-icons/fa6";

export default function Dashboard() {
  const { isLoading } = useAuth()
  const [nodes, setNodes] = useState<Node[]>([])
  const [metadata, setMetadata] = useState<Metadata[]>([])
  const [experiments, setExperiments] = useState<Experiment[]>([])

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      const [fetchedNodes, fetchedMetadata, fetchedExperiments] = await Promise.all([
        nodeService.getAll(),
        metadataService.getAll(),
        experimentService.getAll()
      ])
      setNodes(fetchedNodes)
      setMetadata(fetchedMetadata)
      setExperiments(fetchedExperiments)
    } catch (error) {
      console.error('Failed to fetch data:', error)
    }
  }

  if (isLoading) {
    return <div>Loading...</div>
  }

  return (
    <div className="space-y-8">
      <h1 className="text-3xl md:text-4xl font-bold mb-4 text-center">ICFL</h1>
      <p className="text-base md:text-lg mb-8 text-center">Select an option from the sidebar to get started.</p>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-blue-200 p-6 md:p-8 rounded-lg shadow flex items-center md:h-48">
          <FaCircleNodes className="text-blue-500 text-3xl md:text-5xl mr-4" />
          <div>
            <h2 className="text-xl md:text-2xl font-semibold mb-2">Nodes</h2>
            <p className="text-blue-700 md:text-lg">Total Nodes: {nodes.length}</p>
            <p className="text-blue-700 md:text-lg">Approved Nodes: {nodes.filter(node => node.Approved).length}</p>
            <p className="text-blue-700 md:text-lg">Pending Nodes: {nodes.filter(node => !node.Approved).length}</p>
          </div>
        </div>
        <div className="bg-green-200 p-6 md:p-8 rounded-lg shadow flex items-center md:h-48">
          <FaDatabase className="text-green-500 text-3xl md:text-5xl mr-4" />
          <div>
            <h2 className="text-xl md:text-2xl font-semibold mb-2">Metadata</h2>
            <p className="text-green-700 md:text-lg">Total Metadata Entries: {metadata.length}</p>
            <p className="text-green-700 md:text-lg">Unique Types: {new Set(metadata.map(item => item.Type)).size}</p>
            <p className="text-green-700 md:text-lg">Unique Tags: {new Set(metadata.flatMap(item => item.Tags)).size}</p>
          </div>
        </div>
        <div className="bg-yellow-200 p-6 md:p-8 rounded-lg shadow flex items-center md:h-48">
          <FaFlask className="text-yellow-500 text-3xl md:text-5xl mr-4" />
          <div>
            <h2 className="text-xl md:text-2xl font-semibold mb-2">Experiments</h2>
            <p className="text-yellow-700 md:text-lg">Total Experiments: {experiments.length}</p>
            <p className="text-yellow-700 md:text-lg">Active Experiments: {experiments.filter(exp => exp.Status === 'TRAINING').length}</p>
            <p className="text-yellow-700 md:text-lg">Completed Experiments: {experiments.filter(exp => exp.Status === 'STOPPED').length}</p>
          </div>
        </div>
      </div>
    </div>
  )
}
