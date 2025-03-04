const BASE_URL = process.env.NEXT_PUBLIC_API + '/api';
const BASE_URL_ENTRY = process.env.NEXT_PUBLIC_API;

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
  }
}

async function handleResponse(response: Response) {
  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new ApiError(response.status, error.message || 'An error occurred');
  }
  return response.json();
}

export async function apiCall(endpoint: string, options: RequestInit = {}) {
  const token = localStorage.getItem('authToken');
  const headers = new Headers(options.headers);
  
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  return handleResponse(response);
}

export async function apiCallEntry(endpoint: string, options: RequestInit = {}) {
  const response = await fetch(`${BASE_URL_ENTRY}${endpoint}`, {
    ...options,
  });

  return handleResponse(response);
}