import { useAuth } from '@/context/AuthContext';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import MessageDisplay from './MessageDisplay';

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

interface ReservationFormProps {
  onCancel: () => void;
  onSuccess: (reservation: Reservation) => void;
}

export default function CreateReservation({
  onCancel,
  onSuccess,
}: ReservationFormProps) {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const [games, setGames] = useState<Game[]>([]);
  const { userId, jwt, error, setError } = useAuth();
  const [newReservation, setNewReservation] = useState<NewReservation>({
    Date: '',
    StartTime: '',
    EndTime: '',
    UserID: userId,
    GameID: '',
  });

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const gamesResponse = await axios.get(`${apiUrl}/api/games`);
        setGames(gamesResponse.data);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };
    fetchUserData();
  }, [userId, apiUrl]);

  // const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

  const createReservation = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const response = await axios.post(
        `${apiUrl}/api/reservations`,
        newReservation,
        {
          headers: {
            Authorization: `Bearer ${jwt}`,
          },
        }
      );
      const game = games.find((g) => g.ID === newReservation.GameID);
      const updatedReservation = {
        ...response.data,
        GameName: game ? game.Title : 'Unknown',
      };
      // setReservations([updatedReservation, ...(reservations || [])]);
      setNewReservation({
        Date: '',
        StartTime: '',
        EndTime: '',
        UserID: userId,
        GameID: '',
      });
      onSuccess(updatedReservation);
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        console.error('API Error:', error.response.data.error);
        setError(error.response.data.error);
      } else {
        console.error('Error creating reservation', error);
      }
    }
  };

  return (
    <form
      className='mb-6 p-4 bg-nook-olive rounded shadow flex flex-col w-full'
      onSubmit={createReservation}
    >
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
        {games?.map((g) => (
          <option key={g.ID} value={g.ID}>
            {g.Title}
          </option>
        ))}
      </select>
      <MessageDisplay errorMsg={error} />
      <button
        className='bg-nook-light-olive text-nook-charcoal hover:text-white rounded mb-2 p-2 w-full '
        type='submit'
      >
        Create Reservation
      </button>
      <button
        onClick={onCancel}
        className='bg-nook-dark-rose hover:bg-nook-rose rounded p-2 w-full'
        type='button'
      >
        Cancel
      </button>
    </form>
  );
}
