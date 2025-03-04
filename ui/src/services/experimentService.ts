import { apiCall } from '@/services/api'
import { Node } from '@/services/nodeService'
import { Metadata } from '@/services/metadataService'

export type ExperimentNodeStatus = 
  | "PENDING"
  | "ACCEPTED"
  | "REJECTED"
  | "TRAINING"
  | "STOPPED";

export const ExperimentNodeStatus = {
  PENDING: "PENDING" as ExperimentNodeStatus,
  ACCEPTED: "ACCEPTED" as ExperimentNodeStatus,
  REJECTED: "REJECTED" as ExperimentNodeStatus,
  TRAINING: "TRAINING" as ExperimentNodeStatus,
  STOPPED: "STOPPED" as ExperimentNodeStatus,
};

export interface Experiment {
  ID: number;
  Username: string;
  Name: string;
  Description: string;
  BasePath: string;
  Status: string;
  CreatedAt: string;
  UpdatedAt: string;
  ExperimentNodes: ExperimentNode[];
  User: {
    Username: string
  }
}

export interface ExperimentNode {
  ExperimentID: number;
  NodeID: number;
  MetadataID: number;
  Status: ExperimentNodeStatus;
  Node: Node;
  Metadata: Metadata;
}

export const experimentService = {
  getAll: () => apiCall('/experiments', { method: 'GET' }) as Promise<Experiment[]>,
  
  create: (formData: FormData) => apiCall('/experiments', {
    method: 'POST',
    body: formData,
  }) as Promise<Experiment>,
  
  update: (id: number, formData: FormData) => apiCall(`/experiments/${id}`, {
    method: 'PUT',
    body: formData,
  }) as Promise<Experiment>,

  resendFiles: (id: number, formData: FormData) => apiCall(`/experiments/${id}/update-files`, {
    method: 'POST',
    body: formData,
  }) as Promise<{ status: string }>,

  startTraining: (id: number) => apiCall(`/experiments/${id}/start`, {
    method: 'POST',
  }) as Promise<Experiment>,

  stopTraining: (id: number) => apiCall(`/experiments/${id}/stop`, {
    method: 'POST',
  }) as Promise<Experiment>,
}
