'use client';
import axios from 'axios';
import React, { createContext, useContext, useEffect, useState } from 'react';

interface AuthContextType {
  isAuthenticated: boolean;
  setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>;
  email: string;
  setEmail: React.Dispatch<React.SetStateAction<string>>;
  name: string;
  setName: React.Dispatch<React.SetStateAction<string>>;
  jwt: string;
  setJwt: React.Dispatch<React.SetStateAction<string>>;
  isAdmin: boolean;
  setIsAdmin: React.Dispatch<React.SetStateAction<boolean>>;
  userId: string;
  login: (user: { Email: string; Password: string }) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [email, setEmail] = useState('');
  const [jwt, setJwt] = useState('');
  const [name, setName] = useState('');
  const [isAdmin, setIsAdmin] = useState(false);
  const [userId, setUserId] = useState('');

  interface LoginResponse {
    ID: string;
    Email: string;
    FirstName: string;
    Token: string;
    RefreshToken: string;
    IsAdmin: boolean;
  }

  interface RefreshResponse {
    Token: string;
    Name: string;
    IsAdmin: boolean;
    Email: string;
    ID: string;
  }

  useEffect(() => {
    const attemptAutoLogin = async () => {
      try {
        if (jwt) {
          console.log(jwt);
          setIsAuthenticated(true);
          return;
        }
        console.log('running a refresh req');
        // check if refresh_token cookie exists
        const response = await axios.post<RefreshResponse>(
          `${apiUrl}/api/refresh`,
          {},
          { withCredentials: true }
        );
        // get user data
        setJwt(response.data.Token);
        setIsAuthenticated(true);
        setIsAdmin(response.data.IsAdmin);
        setName(response.data.Name);
        setEmail(response.data.Email);
        setUserId(response.data.ID);
      } catch (error) {
        console.log('not authenticated:', error);
      }
    };
    attemptAutoLogin();
  }, [jwt, apiUrl]);

  const login = async (user: { Email: string; Password: string }) => {
    try {
      const response = await axios.post<LoginResponse>(
        `${apiUrl}/api/users/login`,
        user,
        { withCredentials: true }
      );
      setIsAuthenticated(true);
      setEmail(user.Email);
      setName(response.data.FirstName);
      setJwt(response.data.Token);
      setIsAdmin(response.data.IsAdmin);
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
    try {
      await axios.post(
        `${apiUrl}/api/users/logout`,
        {},
        { withCredentials: true }
      );
      setIsAuthenticated(false);
      setEmail('');
      setJwt('');
      setName('');
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        // setErrorMsg(error.response.data.error);
        console.error('API Error:', error.response.data.error);
      } else {
        // setErrorMsg('Could not login');
        console.error('Error logging out user', error);
      }
    }
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        setIsAuthenticated,
        email,
        setEmail,
        name,
        setName,
        isAdmin,
        setIsAdmin,
        jwt,
        setJwt,
        login,
        logout,
        userId,
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
