import axios from 'axios';
import React, { useEffect, useState } from 'react';

const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function UserData() {
  const [userCount, setUserCount] = useState<number>(0);
  const [usersCountRes, setUsersCountRes] = useState<number>(0);

  useEffect(() => {
    const getUserData = async () => {
      try {
        const response = await axios.get(`${apiUrl}/api/users/stats`);
        setUserCount(response.data?.TotalUsers);
        setUsersCountRes(response.data?.UsersWithReservations);
      } catch (error) {
        console.log(error);
      }
    };
    getUserData();
  });

  return (
    <div className='flex gap-1 mb-1'>
      <p className='bg-nook-charcoal rounded p-4'>
        Registered Users:{' '}
        <span className='font-bold text-nook-light-olive'>{userCount}</span>
      </p>
      <p className='bg-nook-charcoal rounded p-4'>
        Users with at least 1 reservation:{' '}
        <span className='font-bold text-nook-light-olive'>{usersCountRes}</span>
      </p>
    </div>
  );
}
