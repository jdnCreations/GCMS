import { useAuth } from '@/context/AuthContext';
import { useRouter } from 'next/navigation';
import React, { useEffect, useState } from 'react';
import UserBar from './UserBar';
import GameForm from './GameForm';
import GenreForm from './GenreForm';
import axios, { isAxiosError } from 'axios';
import UserData from './UserData';
import ManageUsers from './ManageUsers';

interface PGTime {
  Microseconds: number;
  Valid: boolean;
}

interface Reservation {
  ID: string;
  ResDate: string;
  StartTime: PGTime;
  EndTime: PGTime;
  UserID: string;
  GameID: string;
  Active: boolean;
  Title: string;
}

const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function AdminDashboard() {
  const { isAdmin, isAuthenticated, isLoading } = useAuth();
  const [isManagingUsers, setIsManagingUsers] = useState(false);
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const router = useRouter();

  useEffect(() => {
    let isMounted = true;
    if (!isLoading && !isAuthenticated) {
      router.push('/');
    } else if (!isLoading && isAuthenticated) {
      const getReservations = async () => {
        try {
          const response = await axios.get<Reservation[]>(
            `${apiUrl}/api/reservations/today`
          );
          if (isMounted) {
            console.log(response.data);
            setReservations(response.data);
          }
        } catch (error) {
          if (isAxiosError(error) && error.response?.data?.error) {
            console.log(error.response.data.error);
            setReservations([]);
          } else {
            console.error('unexpected error:', error);
          }
        }
      };
      getReservations();
    }

    return () => {
      isMounted = false;
    };
  }, [isAuthenticated, isLoading, router]);

  const handleCancel = () => {
    setIsManagingUsers(false);
  };

  const formatTime = (microseconds: number) => {
    const milli = microseconds / 1000;
    const date = new Date(milli);
    const hours = date.getUTCHours();
    const minutes = date.getUTCMinutes();

    return new Date(date.setHours(hours, minutes)).toLocaleTimeString('en-AU', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: true,
    });
  };

  if (!isAuthenticated || !isAdmin) {
    return <div>Loading...</div>;
  }

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <div className='w-full'>
      <UserBar />
      <div className='bg-nook-light-charcoal rounded-b p-4 mb-1'>
        <p>Reservation Data</p>
        <p>{reservations?.length} reservations today</p>
        {reservations
          ?.slice()
          .sort((a, b) => a.StartTime.Microseconds - b.StartTime.Microseconds)
          .map((reservation) => (
            <div
              className='bg-nook-charcoal p-1 px-2 rounded-sm mb-1'
              key={reservation.ID}
            >
              {reservation.Title} @{' '}
              {formatTime(reservation.StartTime.Microseconds)}
            </div>
          ))}
      </div>
      {!isManagingUsers && (
        <button
          className='bg-nook-charcoal hover:bg-nook-light-charcoal w-full rounded mb-1 p-2'
          onClick={() => setIsManagingUsers(true)}
        >
          Manage Users
        </button>
      )}
      {isManagingUsers && <ManageUsers onCancel={handleCancel} />}
      <UserData />
      <GameForm />
      <GenreForm />
    </div>
  );
}
