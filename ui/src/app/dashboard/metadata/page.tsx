'use client'

import { useState, useEffect } from 'react'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { metadataService, Metadata } from '@/services/metadataService'

type SortKey = 'NodeID' | 'Name' | 'Type' | 'Tags';

export default function MetadataPage() {
  const [metadata, setMetadata] = useState<Metadata[]>([])
  const [sortKey, setSortKey] = useState<SortKey>('NodeID')
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc')
  const [searchTerm, setSearchTerm] = useState('')
  const [typeFilter, setTypeFilter] = useState<string>('all')
  const [tagFilter, setTagFilter] = useState<string>('all')

  useEffect(() => {
    fetchMetadata()
  }, [])

  const fetchMetadata = async () => {
    try {
      const fetchedMetadata = await metadataService.getAll()
      setMetadata(fetchedMetadata)
    } catch (error) {
      console.error('Failed to fetch metadata:', error)
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

  const filteredAndSortedMetadata = metadata
    .filter(item => 
      (searchTerm === '' || 
        Object.values(item).some(value => 
          value.toString().toLowerCase().includes(searchTerm.toLowerCase())
        )
      ) &&
      (typeFilter === 'all' || item.Type === typeFilter) &&
      (tagFilter === 'all' || item.Tags.includes(tagFilter))
    )
    .sort((a, b) => {
      let compareA = a[sortKey];
      let compareB = b[sortKey];

      if (sortKey === 'Tags') {
        compareA = a.Tags.join(', ');
        compareB = b.Tags.join(', ');
      }

      if (compareA < compareB) return sortDirection === 'asc' ? -1 : 1;
      if (compareA > compareB) return sortDirection === 'asc' ? 1 : -1;
      return 0;
    });

  const uniqueTypes = ['all', ...Array.from(new Set(metadata.map(item => item.Type)))]
  const uniqueTags = ['all', ...Array.from(new Set(metadata.flatMap(item => item.Tags)))]

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Metadata</h1>
      <div className="flex space-x-4">
        <Input
          placeholder="Search..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="max-w-sm"
        />
        <Select onValueChange={setTypeFilter} value={typeFilter}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Filter by type" />
          </SelectTrigger>
          <SelectContent>
            {uniqueTypes.map(type => (
              <SelectItem key={type} value={type}>{type}</SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Select onValueChange={setTagFilter} value={tagFilter}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Filter by tag" />
          </SelectTrigger>
          <SelectContent>
            {uniqueTags.map(tag => (
              <SelectItem key={tag} value={tag}>{tag}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>
              <Button variant="ghost" onClick={() => handleSort('NodeID')}>
                Node ID {sortKey === 'NodeID' && (sortDirection === 'asc' ? '▲' : '▼')}
              </Button>
            </TableHead>
            <TableHead>Node Metadata ID</TableHead>
            <TableHead>
              <Button variant="ghost" onClick={() => handleSort('Name')}>
                Name {sortKey === 'Name' && (sortDirection === 'asc' ? '▲' : '▼')}
              </Button>
            </TableHead>
            <TableHead>
              <Button variant="ghost" onClick={() => handleSort('Type')}>
                Type {sortKey === 'Type' && (sortDirection === 'asc' ? '▲' : '▼')}
              </Button>
            </TableHead>
            <TableHead>
              <Button variant="ghost" onClick={() => handleSort('Tags')}>
                Tags {sortKey === 'Tags' && (sortDirection === 'asc' ? '▲' : '▼')}
              </Button>
            </TableHead>
            <TableHead>Description</TableHead>
            <TableHead>Extras</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {filteredAndSortedMetadata.map((item) => (
            <TableRow key={item.NodeMetadataID}>
              <TableCell>{item.NodeID}</TableCell>
              <TableCell>{item.NodeMetadataID}</TableCell>
              <TableCell>{item.Name}</TableCell>
              <TableCell>{item.Type}</TableCell>
              <TableCell>{item.Tags}</TableCell>
              <TableCell>{item.Description}</TableCell>
              <TableCell>{JSON.stringify(item.Extras)}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
