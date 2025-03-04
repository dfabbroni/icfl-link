import { apiCall } from '@/services/api'

export interface Metadata {
  ID: string;
  NodeID: string;
  NodeMetadataID: string;
  Name: string;
  Type: string;
  Tags: string[];
  Description: string;
  Extras: Record<string, any>;
}

export const metadataService = {
  getAll: () => apiCall('/metadata', { method: 'GET' }) as Promise<Metadata[]>,
}
