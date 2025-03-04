'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Checkbox } from "@/components/ui/checkbox"
import { experimentService } from '@/services/experimentService'
import { metadataService, Metadata } from '@/services/metadataService'

interface Node {
  id: string;
  metadata: Metadata[];
}

interface SelectedNode {
    id: string;
    node_id: string;
    metadata_id: string;
  }

export default function CreateExperimentPage() {
  const router = useRouter()
  const [step, setStep] = useState(1)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [zipFile, setZipFile] = useState<File | null>(null)
  const [nodes, setNodes] = useState<Node[]>([])
  const [selectedNodes, setSelectedNodes] = useState<SelectedNode[]>([]);
  const [filter, setFilter] = useState('')

  useEffect(() => {
    fetchMetadata()
  }, [])

  const fetchMetadata = async () => {
    try {
      const metadata = await metadataService.getAll()
      const groupedNodes: {[key: string]: Metadata[]} = {}
      metadata.forEach((item: Metadata) => {
        if (!groupedNodes[item.NodeID]) {
          groupedNodes[item.NodeID] = []
        }
        groupedNodes[item.NodeID].push(item)
      })
      setNodes(Object.entries(groupedNodes).map(([id, metadata]) => ({ id, metadata })))
    } catch (error) {
      console.error('Failed to fetch metadata:', error)
    }
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setZipFile(e.target.files[0])
    }
  }

  const handleNodeSelection = (nodeId: string, id: string, metadataId: string) => {
    setSelectedNodes(prev => {
      const existingIndex = prev.findIndex(node => node.node_id === nodeId);
      if (existingIndex !== -1) {
        const newNodes = [...prev];
        newNodes[existingIndex] = { node_id: nodeId, id: id, metadata_id: metadataId };
        return newNodes;
      } else {
        return [...prev, { node_id: nodeId, id: id, metadata_id: metadataId }];
      }
    });
  }

  const filteredNodes = nodes.filter(node => 
    node.id.toLowerCase().includes(filter.toLowerCase()) ||
    node.metadata.some(m => m.Name.toLowerCase().includes(filter.toLowerCase()))
  )

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const formData = new FormData()
    formData.append('name', name)
    formData.append('description', description)
    if (zipFile) formData.append('experimentFiles', zipFile)
    formData.append('selectedNodes', JSON.stringify(selectedNodes.map(node => ({
        ...node,
        node_id: Number(node.node_id),
        id: Number(node.id),
        metadata_id: Number(node.metadata_id)
    }))));
    
    // Fix the linter error by using Array.from()
    console.log('FormData contents:');
    Array.from(formData.entries()).forEach(([key, value]) => {
      console.log(key, typeof value === 'string' ? value : `File: ${value.name}`);
    });

    try {
      await experimentService.create(formData)
      router.push('/dashboard/experiments')
    } catch (error) {
      console.error('Failed to create experiment:', error)
    }
  }

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Create New Experiment</h1>
      {step === 1 && (
        <form onSubmit={(e) => { e.preventDefault(); setStep(2) }} className="space-y-4">
          <Input 
            placeholder="Experiment Name" 
            value={name} 
            onChange={(e) => setName(e.target.value)} 
            required 
          />
          <Textarea 
            placeholder="Description" 
            value={description} 
            onChange={(e) => setDescription(e.target.value)} 
            required 
          />
          <Input 
            type="file" 
            onChange={handleFileChange} 
            accept=".zip"
            required
          />
          <Button type="submit">Next</Button>
        </form>
      )}
      {step === 2 && (
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input 
            placeholder="Filter nodes" 
            value={filter} 
            onChange={(e) => setFilter(e.target.value)} 
          />
          <div className="space-y-2">
            {filteredNodes.map(node => (
              <div key={node.id} className="border p-2 rounded">
                <h3 className="font-bold">{node.id}</h3>
                {node.metadata.map(metadata => (
                  <div key={metadata.ID} className="flex items-center">
                    <Checkbox
                      id={metadata.ID}
                      checked={selectedNodes.some(n => n.node_id === node.id && n.id === metadata.ID)}
                      onCheckedChange={() => handleNodeSelection(node.id, metadata.ID, metadata.NodeMetadataID)}
                    />
                    <label htmlFor={metadata.NodeMetadataID} className="ml-2">
                      {metadata.Name} ({metadata.Type})
                    </label>
                  </div>
                ))}
              </div>
            ))}
          </div>
          <Button type="submit">Create Experiment</Button>
        </form>
      )}
    </div>
  )
}
