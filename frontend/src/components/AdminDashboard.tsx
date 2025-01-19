import { useAuth } from '@/context/AuthContext';
import { useRouter } from 'next/navigation';
import React, { useEffect } from 'react';
import UserBar from './UserBar';

export default function AdminDashboard() {
  const { isAdmin, isAuthenticated, isLoading, logout } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/');
    }
  }, [isAuthenticated, isLoading, router]);

  if (!isAuthenticated || !isAdmin) {
    return <div>Loading...</div>;
  }

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <div className='w-full'>
      <UserBar />
      <p>Reservation Data</p>
      <div className='bg-white text-gray-800 font-bold rounded'>
        <div>
          <p>Active Reservations</p>
          <p>12</p>
        </div>
        <div>
          <p>Registered Users</p>
          <p>3</p>
        </div>
        <div>
          <p>Users who have made reservation(s)</p>
          <p>1</p>
        </div>
      </div>
    </div>
  );
}
