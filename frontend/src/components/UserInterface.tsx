'use client';

import React, { useEffect, useState } from 'react';
import CardComponent from './CardComponent';
import axios from 'axios';
import Login from './Login';
import { useAuth } from '@/context/AuthContext';
import CreateUser from './Register';
import UpdateUser from './UpdateUser';
import { useRouter } from 'next/navigation';

interface User {
  ID: string;
  FirstName: string;
  LastName: string;
  Email: string;
  Password: string;
}

interface CurrentUser {
  ID: string;
  Email: string;
  Token: string;
}

const UserInterface: React.FC = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const { isAuthenticated, name, jwt, email, logout } = useAuth();
  const router = useRouter();

  const handleLogout = async () => {
    await logout();
  };

  return (
    <div className=' justify-center items-center flex flex-col rounded-lg py-2 px-8'>
      <p>hello {name}</p>
      <button onClick={handleLogout}>Logout</button>
      <button onClick={() => router.push('/dashboard')}>Go to dashboard</button>
    </div>
  );
};

export default UserInterface;
