'use client'

import { useState, useEffect, useRef } from 'react'
import Link from 'next/link'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { experimentService, Experiment, ExperimentNode } from '@/services/experimentService'

type SortKey = 'ID' | 'Username';

export default function ExperimentsPage() {
  const [experiments, setExperiments] = useState<Experiment[]>([])
  const [selectedExperiment, setSelectedExperiment] = useState<Experiment | null>(null)
  const [expandedExperimentId, setExpandedExperimentId] = useState<number | null>(null)
  const [isUpdateDialogOpen, setIsUpdateDialogOpen] = useState(false)
  const [sortKey, setSortKey] = useState<SortKey>('ID')
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc')

  const initialModelRef = useRef<HTMLInputElement>(null)
  const taskRef = useRef<HTMLInputElement>(null)
  const clientAppRef = useRef<HTMLInputElement>(null)
  const pyprojectTomlRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    fetchExperiments()
  }, [])

  const fetchExperiments = async () => {
    try {
      const fetchedExperiments = await experimentService.getAll()
      setExperiments(fetchedExperiments)
    } catch (error) {
      console.error('Failed to fetch experiments:', error)
    }
  }

  const handleRowClick = (experiment: Experiment) => {
    if (expandedExperimentId === experiment.ID) {
      setExpandedExperimentId(null)
    } else {
      setExpandedExperimentId(experiment.ID)
    }
  }

  const handleUpdate = (experiment: Experiment, e: React.MouseEvent) => {
    e.stopPropagation()
    setSelectedExperiment(experiment)
    setIsUpdateDialogOpen(true)
  }

  const handleUpdateSubmit = async (e: React.FormEvent) => {
    console.log('handleUpdateSubmit')
    e.preventDefault()
    if (!selectedExperiment) return

    const formData = new FormData()

    if (initialModelRef.current?.files?.[0]) {
      formData.append('initialModel', initialModelRef.current.files[0])
    }
    if (clientAppRef.current?.files?.[0]) {
      formData.append('clientApp', clientAppRef.current.files[0])
    }
    if (pyprojectTomlRef.current?.files?.[0]) {
      formData.append('pyprojectToml', pyprojectTomlRef.current.files[0])
    }
    if (taskRef.current?.files?.[0]) {
      formData.append('taskRef', taskRef.current.files[0])
    }

    try {
      await experimentService.update(selectedExperiment.ID, formData)
      await fetchExperiments()
      setIsUpdateDialogOpen(false)
    } catch (error) {
      console.error('Failed to update experiment:', error)
    }
  }

  const handleResendFiles = async (experimentId: number, e: React.MouseEvent) => {
    e.stopPropagation()
    const formData = new FormData()
    await experimentService.resendFiles(experimentId, formData)
  }

  const handleStartTraining = async (experimentId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      const updatedExperiment = await experimentService.startTraining(experimentId);
      setExperiments(prevExperiments => 
        prevExperiments.map(exp => exp.ID === experimentId ? updatedExperiment : exp)
      );

      // Ensure the expanded state remains the same
      setExpandedExperimentId(prevId => (prevId === experimentId ? experimentId : prevId));
    } catch (error) {
      console.error('Failed to start training:', error);
    }
  }

  const handleStopTraining = async (experimentId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      const updatedExperiment = await experimentService.stopTraining(experimentId);
      setExperiments(prevExperiments => 
        prevExperiments.map(exp => exp.ID === experimentId ? updatedExperiment : exp)
      );

      // Ensure the expanded state remains the same
      setExpandedExperimentId(prevId => (prevId === experimentId ? experimentId : prevId));
    } catch (error) {
      console.error('Failed to stop training:', error);
    }
  }

  const handleSort = (key: SortKey) => {
    if (key === sortKey) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')
    } else {
      setSortKey(key)
      setSortDirection('asc')
    }
  }

  const sortedExperiments = [...experiments].sort((a, b) => {
    if (a[sortKey] < b[sortKey]) return sortDirection === 'asc' ? -1 : 1
    if (a[sortKey] > b[sortKey]) return sortDirection === 'asc' ? 1 : -1
    return 0
  })

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Experiments</h1>
        <Link href="/dashboard/experiments/create">
          <Button>Create New Experiment</Button>
        </Link>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>
              <Button variant="ghost" onClick={() => handleSort('ID')}>
                ID {sortKey === 'ID' && (sortDirection === 'asc' ? '▲' : '▼')}
              </Button>
            </TableHead>
            <TableHead>
              <Button variant="ghost" onClick={() => handleSort('Username')}>
                Username {sortKey === 'Username' && (sortDirection === 'asc' ? '▲' : '▼')}
              </Button>
            </TableHead>
            <TableHead>Name</TableHead>
            <TableHead>Description</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Created At</TableHead>
            <TableHead>Updated At</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {sortedExperiments.map((experiment) => (
            <>
              <TableRow 
                key={experiment.ID} 
                onClick={() => handleRowClick(experiment)}
                className="cursor-pointer hover:bg-gray-100"
              >
                <TableCell>{experiment.ID}</TableCell>
                <TableCell>{experiment.User.Username}</TableCell>
                <TableCell>{experiment.Name}</TableCell>
                <TableCell>{experiment.Description}</TableCell>
                <TableCell>{experiment.Status}</TableCell>
                <TableCell>{new Date(experiment.CreatedAt).toLocaleString()}</TableCell>
                <TableCell>{new Date(experiment.UpdatedAt).toLocaleString()}</TableCell>
              </TableRow>
              {expandedExperimentId === experiment.ID && (
                <TableRow>
                  <TableCell colSpan={7}>
                    <div className="p-4 bg-gray-50">
                      <h3 className="font-bold mb-2">Experiment Nodes:</h3>
                      {experiment.ExperimentNodes && experiment.ExperimentNodes.map((node: ExperimentNode) => (
                        <div key={`${node.ExperimentID}-${node.NodeID}`} className="mb-2">
                          <p>Node ID: {node.NodeID}</p>
                          <p>Node: {node.Node.Username}</p>
                          <p>Metadata: {node.Metadata.Name}</p>
                          <p>Status: {node.Status}</p>
                        </div>
                      ))}
                      <div className="mt-4 space-x-2">
                        <Button onClick={(e) => handleUpdate(experiment, e)}>Update Experiment</Button>
                        <Button onClick={(e) => handleResendFiles(experiment.ID, e)}>Resend Files</Button>
                        <Button onClick={(e) => handleStartTraining(experiment.ID, e)}>Start Training</Button>
                        <Button onClick={(e) => handleStopTraining(experiment.ID, e)}>Stop Training</Button>
                      </div>
                    </div>
                  </TableCell>
                </TableRow>
              )}
            </>
          ))}
        </TableBody>
      </Table>

      <Dialog open={isUpdateDialogOpen} onOpenChange={setIsUpdateDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Update Experiment</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleUpdateSubmit} className="space-y-4">
            <div>
              <label htmlFor="initialModel">Initial Model:</label>
              <Input id="initialModel" type="file" ref={initialModelRef} accept=".pt" />
            </div>
            <div>
              <label htmlFor="clientApp">Client App:</label>
              <Input id="clientApp" type="file" ref={clientAppRef} accept=".py" />
            </div>
            <div>
              <label htmlFor="pyprojectToml">Pyproject TOML:</label>
              <Input id="pyprojectToml" type="file" ref={pyprojectTomlRef} accept=".toml" />
            </div>
            <div>
              <label htmlFor="taskRef">Task:</label>
              <Input id="taskRef" type="file" ref={taskRef} accept=".py" />
            </div>
            <Button type="submit">Update Experiment</Button>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
