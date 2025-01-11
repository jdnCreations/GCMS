'use client';
import ReservationComponent from '@/components/ReservationComponent';
import axios from 'axios';
import React, { use, useEffect, useState } from 'react';

interface User {
  ID: string;
  FirstName: string;
  LastName: string;
  Email: string;
}

interface NewReservation {
  ID?: string;
  Date: string;
  StartTime: string;
  EndTime: string;
  UserID: string;
  GameID: string;
  GameName?: string;
}

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

interface Game {
  ID: string;
  Title: string;
  Copies: number;
}

function Page({ params }: { params: Promise<{ id: string }> }) {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const { id } = use(params);
  const [user, setUser] = useState<User>();
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [games, setGames] = useState<Game[]>([]);
  const [newReservation, setNewReservation] = useState<NewReservation>({
    Date: '',
    StartTime: '',
    EndTime: '',
    UserID: id,
    GameID: '',
  });

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const [userResponse, reservationsResponse, gamesResponse] =
          await Promise.all([
            axios.get(`${apiUrl}/api/users/${id}`),
            axios.get(`${apiUrl}/api/reservations/${id}`),
            axios.get(`${apiUrl}/api/games`),
          ]);
        setUser(userResponse.data);
        setReservations(reservationsResponse.data);
        setGames(gamesResponse.data);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };
    fetchUserData();
  }, [id, apiUrl]);

  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  console.log(timezone); // Example: "America/New_York"

  const createReservation = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log(newReservation);
    try {
      const response = await axios.post(
        `${apiUrl}/api/reservations`,
        newReservation
      );
      const game = games.find((g) => g.ID === newReservation.GameID);
      const updatedReservation = {
        ...response.data,
        GameName: game ? game.Title : 'Unknown',
      };
      setReservations([updatedReservation, ...(reservations || [])]);
      setNewReservation({
        Date: '',
        StartTime: '',
        EndTime: '',
        UserID: id,
        GameID: '',
      });
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        console.error('API Error:', error.response.data.error);
      } else {
        console.error('Error creating reservation', error);
      }
    }
  };

  const deleteReservation = async (id: string) => {
    try {
      await axios.delete(`${apiUrl}/api/reservations/${id}`);
      setReservations(reservations.filter((res) => res.ID !== id));
    } catch (error) {
      console.error('Error deleting reservation:', error);
    }
  };

  return (
    <div className='flex flex-col items-center gap-4'>
      <div>
        <h1>User Info</h1>
        <p>
          {user?.FirstName} {user?.LastName}
        </p>
      </div>
      <div className='flex flex-col'>
        <form
          className='mb-6 p-4 bg-blue-400 rounded shadow flex flex-col'
          onSubmit={createReservation}
        >
          <h2 className='font-bold pb-2'>
            Create a reservation for {user?.FirstName}
          </h2>
          <label htmlFor='date'>Reservation Date</label>
          <input
            name='date'
            type='date'
            value={newReservation?.Date}
            onChange={(e) =>
              setNewReservation({
                ...newReservation,
                Date: e.target.value,
              })
            }
            className='mb-2 p-2 border border-gra-300 rounded text-gray-800'
          />
          <label htmlFor='startTime'>Start time</label>
          <input
            name='startTime'
            type='time'
            value={newReservation?.StartTime}
            onChange={(e) =>
              setNewReservation({
                ...newReservation,
                StartTime: e.target.value,
              })
            }
            className='mb-2 p-2 border border-gray-300 rounded text-gray-800'
          />

          <label htmlFor='endTime'>End time</label>
          <input
            type='time'
            value={newReservation?.EndTime}
            onChange={(e) =>
              setNewReservation({
                ...newReservation,
                EndTime: e.target.value,
              })
            }
            className='mb-2 p-2 border border-gray-300 roudned text-gray-800'
          />
          <label htmlFor='game'>Game</label>
          <select
            className='mb-2 p-2 border text-gray-900 border-gray-300 rounded'
            name='game'
            id='game'
            onChange={(e) => {
              e.preventDefault();
              const selectedGameId = e.target.value;
              setNewReservation({ ...newReservation, GameID: selectedGameId });
            }}
          >
            <option value=''>Select a Game</option>
            {games.map((g) => (
              <option key={g.ID} value={g.ID}>
                {g.Title}
              </option>
            ))}
          </select>
          <button type='submit'>Create Reservation</button>
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

export default Page;
