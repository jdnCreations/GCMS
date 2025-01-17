'use client';
import ReservationComponent from '@/components/ReservationComponent';
import { useAuth } from '@/context/AuthContext';
import axios, { AxiosError } from 'axios';
import { useRouter } from 'next/navigation';
import React, { useEffect, useState } from 'react';

interface Reservation {
  ID?: string;
  ResDate: string;
  StartTime: {
    Microseconds: number;
    Valid: boolean;
  };
  EndTime: {
    Microseconds: number;
    Valid: boolean;
  };
  UserID: string;
  GameID: string;
  GameName?: string;
}

interface UpdatedUser {
  FirstName?: string;
  LastName?: string;
  Email?: string;
}

const MessageDisplay = ({ updateMsg, errorMsg }) => {
  const message = errorMsg || updateMsg || '';
  return (
    <div
      className={`p-2 rounded my-1 ${errorMsg ? 'bg-red-500' : 'bg-green-500'}`}
    >
      {message}
    </div>
  );
};

export default function Dashboard() {
  const { name, isAuthenticated, isAdmin, userId } = useAuth();
  const [updatedUserInfo, setUpdatedUserInfo] = useState<UpdatedUser>({
    FirstName: '',
    LastName: '',
    Email: '',
  });
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const [updateMsg, setUpdateMsg] = useState('');
  const [errMsg, setErrMsg] = useState('');
  const router = useRouter();

  useEffect(() => {
    const getReservations = async () => {
      try {
        const response = await axios.get(
          `${apiUrl}/api/reservations/${userId}`
        );
        setReservations(response.data);
      } catch (error) {
        console.error(error);
      }
    };
    if (!userId) {
      console.log('no userId available');
      return;
    }
    getReservations();
  }, [apiUrl, userId]);

  const deleteReservation = async (id: string) => {
    try {
      await axios.delete(`${apiUrl}/api/reservations/${id}`);
      setReservations(reservations.filter((res) => res.ID !== id));
    } catch (error) {
      console.error('Error deleting reservation:', error);
    }
  };

  const updateUser = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const response = await axios.put(
        `${apiUrl}/api/users/${userId}`,
        updatedUserInfo
      );
      setUpdatedUserInfo({ FirstName: '', LastName: '', Email: '' });
      // tell user their info was updated successfully
      setUpdateMsg('updated successfully');
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        setErrMsg(error.response.data);
      }
      console.log(error);
    }
  };

  return (
    <div className='bg-[#D8CBAF] min-h-screen'>
      <p>Your dashboard {userId}</p>
      <div className='flex flex-col items-center gap-2'>
        <form
          onSubmit={updateUser}
          className='mb-6 p-4 bg-[#B7A99A] rounded shadow text-gray-800'
        >
          <h1 className='font-bold text-2xl'>Update User Details</h1>
          <MessageDisplay errorMsg={errMsg} updateMsg={updateMsg} />
          <input
            type='text'
            placeholder='First Name'
            value={updatedUserInfo.FirstName}
            onChange={(e) =>
              setUpdatedUserInfo({
                ...updatedUserInfo,
                FirstName: e.target.value,
              })
            }
            className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='Last Name'
            value={updatedUserInfo.LastName}
            onChange={(e) =>
              setUpdatedUserInfo({
                ...updatedUserInfo,
                LastName: e.target.value,
              })
            }
            className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='Email'
            value={updatedUserInfo.Email}
            onChange={(e) =>
              setUpdatedUserInfo({ ...updatedUserInfo, Email: e.target.value })
            }
            className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded'
          />
          <button
            className='bg-[#A8BBA1] text-white w-full p-2 rounded hover:bg-[#E4B7B2]'
            type='submit'
          >
            Update Details
          </button>
        </form>
      </div>
      <div className='flex flex-col items-center gap-2'>
        <h1 className='font-bold text-2xl'>Reservations</h1>
        {reservations?.map((reservation) => (
          <ReservationComponent
            key={reservation.ID}
            reservation={reservation}
            onDelete={deleteReservation}
          />
        ))}
      </div>
    </div>
  );
}
