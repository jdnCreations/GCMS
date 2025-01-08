'use client';
import axios from 'axios';
import React, { use, useEffect, useState } from 'react';

interface User {
  ID: string;
  FirstName: string;
  LastName: string;
  Email: string;
}

interface Reservation {
  ID?: string;
  StartTime: string;
  EndTime: string;
  UserID: string;
  GameID: string;
}

function Page({ params }: { params: Promise<{ id: string }> }) {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const { id } = use(params);
  const [user, setUser] = useState<User>();
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [newReservation, setNewReservation] = useState<Reservation>({
    StartTime: '',
    EndTime: '',
    UserID: id,
    GameID: '',
  });

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const [userResponse, reservationsResponse] = await Promise.all([
          axios.get(`${apiUrl}/api/users/${id}`),
          axios.get(`${apiUrl}/api/reservations/${id}`),
        ]);
        setUser(userResponse.data);
        setReservations(reservationsResponse.data);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };
    fetchUserData();
  }, [id, apiUrl]);

  const createReservation = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const response = await axios.post(
        `${apiUrl}/api/reservation`,
        newReservation
      );
      setReservations([response.data, ...reservations]);
      setNewReservation({
        endTime: '',
        startTime: '',
        userID: '',
        gameID: '',
      });
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        console.error('API Error:', error.response.data.error);
      } else {
        console.error('Error creating reservation', error);
      }
    }
  };

  console.log(reservations);

  return (
    <div>
      <div>
        <h1>User Info</h1>
        <p>
          {user?.FirstName} {user?.LastName}
        </p>
      </div>
      <div>
        <form
          className='mb-6 p-4 bg-blue-400 rounded shadow'
          onSubmit={createReservation}
        >
          <input
            type='datetime-local'
            value={newReservation?.start_time}
            onChange={(e) =>
              setNewReservation({
                ...newReservation,
                start_time: e.target.value,
              })
            }
            className='mb-2 p-2 border border-gray-300 roudned'
          />
          <input
            type='datetime-local'
            value={newReservation?.end_time}
            onChange={(e) =>
              setNewReservation({
                ...newReservation,
                end_time: e.target.value,
              })
            }
            className='mb-2 p-2 border border-gray-300 roudned'
          />
          <button type='submit'>Create Reservation</button>
        </form>
      </div>

      <div>
        <h1>Reservations</h1>
        {reservations.map((reservation) => (
          <div key={reservation.ID}>
            <p>{reservation.StartTime}</p>
            <p>{reservation.EndTime}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

export default Page;
