'use client';
import axios from 'axios';
// src/context/AuthContext.tsx

import React, { createContext, useContext, useState } from 'react';

interface AuthContextType {
  isAuthenticated: boolean;
  email: string;
  jwt: string;
  login: (user: { Email: string; Password: string }) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [email, setEmail] = useState('');
  const [jwt, setJwt] = useState('');

  const login = async (user: { Email: string; Password: string }) => {
    try {
      const response = await axios.post(`${apiUrl}/api/users/login`, user);
      setIsAuthenticated(true);
      setEmail(user.Email);
      setJwt(response.data.Token);
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        // setErrorMsg(error.response.data.error);
        console.error('API Error:', error.response.data.error);
      } else {
        // setErrorMsg('Could not login');
        console.error('Error logging in user', error);
      }
    }
  };

  const logout = async () => {
    setIsAuthenticated(false);
    setEmail('');
    setJwt('');
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        email,
        jwt,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
