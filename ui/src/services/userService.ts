import { apiCallEntry } from './api';

export interface UserCredentials {
  Username: string;
  Password: string;
}

export interface AuthResponse {
  token: string | null;
}

export const userService = {
  login: (credentials: UserCredentials) => 
    apiCallEntry('/users/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    }) as Promise<AuthResponse>,

  register: (credentials: UserCredentials) => 
    apiCallEntry('/users', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    }) as Promise<void>,
};
