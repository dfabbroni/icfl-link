import { apiCall, apiCallEntry } from './api';

export interface NodeCredentials {
  Username: string;
  Password: string;
  PublicKey: string;
}

export interface NodeStatusResponse {
  Approved: boolean;
}

export interface Node {
  ID: string;
  Username: string;
  Approved: boolean;
  LastSeen: string;
}

export const nodeService = {
  login: (credentials: Omit<NodeCredentials, 'PublicKey'>) => 
    apiCallEntry('/nodes/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    }) as Promise<NodeStatusResponse>,

  register: (credentials: NodeCredentials) => 
    apiCallEntry('/nodes', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    }) as Promise<NodeStatusResponse>,

  getAll: () =>
    apiCall('/nodes', {
      method: 'GET'
    }) as Promise<Node[]>,
  
  approve: (id: string) => 
    apiCall(`/nodes/${id}/accept`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      }, }) as Promise<Node>,
    
  reject: (id: string) =>
    apiCall(`/nodes/${id}`, {
      method: 'DELETE'
    }) as Promise<void>,
  
};
