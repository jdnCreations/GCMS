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
  isLoading: boolean;
  error: string | null;
  setError: React.Dispatch<React.SetStateAction<string | null>>;
  register: (
    user: {
      FirstName: string;
      LastName: string;
      Email: string;
      Password: string;
    },
    e: React.FormEvent<HTMLFormElement>
  ) => Promise<void>;
  login: (
    user: { Email: string; Password: string },
    e: React.FormEvent<HTMLFormElement>
  ) => Promise<void>;
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
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

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
      setIsLoading(true);
      try {
        if (jwt) {
          // make req to server to validate jwt
          await axios.get(`${apiUrl}/api/verify`, {
            headers: {
              Authorization: `Bearer ${jwt}`,
            },
          });
          setIsAuthenticated(true);
          setIsLoading(false);
          return;
        }
      } catch (error) {
        console.error('no valid access token, attempting refresh', error);
      }
      try {
        // check if refresh_token cookie exists
        console.log('checking refresh token');
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
      } finally {
        setIsLoading(false);
      }
    };
    attemptAutoLogin();
  }, [jwt, apiUrl]);

  const register = async (
    user: {
      FirstName: string;
      LastName: string;
      Email: string;
      Password: string;
    },
    e: React.FormEvent<HTMLFormElement>
  ) => {
    e.preventDefault();

    try {
      await axios.post(`${apiUrl}/api/users`, user);
      await login({ Email: user.Email, Password: user.Password });
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        console.error('API Error:', error.response.data.error);
        setError(error.response.data.error);
      } else {
        console.error('Error creating user:', error);
      }
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (user: { Email: string; Password: string }) => {
    setIsLoading(true);
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
      setUserId(response.data.ID);
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        // setErrorMsg(error.response.data.error);
        console.error('API Error:', error.response.data.error);
        setError(error.response.data.error);
      } else {
        if (error instanceof Error) {
          console.error('Error logging in user', error);
          setError(error.message);
        }
      }
    } finally {
      setIsLoading(false);
      // setError(null);
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
      setUserId('');
      setIsAdmin(false);
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
        register,
        login,
        logout,
        isLoading,
        error,
        setError,
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
