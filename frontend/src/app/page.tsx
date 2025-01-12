'use client';

import UserInterface from '@/components/UserInterface';
import { useAuth } from '@/context/AuthContext';
import { createContext } from 'react';

export const AuthContext = createContext(null);

export default function Home() {
  const { isAuthenticated, login, logout } = useAuth();
  return (
    <div className='grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]'>
      <main className='flex flex-col gap-8 row-start-2 items-center sm:items-start'>
        <h1>{isAuthenticated ? 'Logged In' : 'Logged Out'}</h1>
        {isAuthenticated ? (
          <button onClick={logout}>Logout</button>
        ) : (
          <button onClick={login}>Login</button>
        )}
        <UserInterface />
      </main>
    </div>
  );
}
