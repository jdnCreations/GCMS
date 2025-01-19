import { useAuth } from '@/context/AuthContext';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import ReservationComponent from './ReservationComponent';
import CreateReservation from './CreateReservation';
import UserBar from './UserBar';

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

export default function UserDashboard() {
  const { userId, jwt } = useAuth();

  const [isCreating, setIsCreating] = useState(false);
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

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
      return;
    }
    getReservations();
  }, [apiUrl, userId]);

  const deleteReservation = async (id: string) => {
    try {
      await axios.delete(`${apiUrl}/api/reservations/${id}`, {
        headers: {
          Authorization: `Bearer ${jwt}`,
        },
      });
      setReservations(reservations.filter((res) => res.ID !== id));
    } catch (error) {
      console.error('Error deleting reservation:', error);
    }
  };

  const handleCancel = () => {
    setIsCreating(false);
  };

  const handleSuccess = (reservation: Reservation) => {
    setReservations([...reservations, reservation]);
    setIsCreating(false);
  };

  return (
    <div className='max-w-sm w-full'>
      <UserBar />

      <div className='flex flex-col gap-2 mt-2'>
        <h1 className='text-nook-charcoal text-2xl font-bold'>
          {isCreating ? 'New Reservation' : 'My Reservations'}
        </h1>
        {!isCreating && (
          <button
            className='bg-nook-light-charcoal hover:bg-nook-charcoal p-2 rounded mb-2'
            onClick={() => setIsCreating(true)}
          >
            Make New Reservation
          </button>
        )}
      </div>

      {isCreating ? (
        <CreateReservation onCancel={handleCancel} onSuccess={handleSuccess} />
      ) : (
        <div className='flex flex-col gap-2'>
          {reservations?.length === 0 || reservations == null ? (
            <p className='text-nook-charcoal text-center'>
              No reservations yet. Create your first one!
            </p>
          ) : (
            reservations?.map((reservation) => (
              <ReservationComponent
                key={reservation.ID}
                onDelete={deleteReservation}
                reservation={reservation}
              />
            ))
          )}
        </div>
      )}

      {/* <div>
        <UpdateUser />
      </div>
      <div className='flex flex-col items-center gap-2'>
        <h1 className='font-bold text-2xl text-nook-charcoal'>Reservations</h1>
        {reservations?.length == 0 ||
          (reservations == null && <p>You currently have no reservations.</p>)}
        {reservations?.map((reservation) => (
          <ReservationComponent
            key={reservation.ID}
            reservation={reservation}
            onDelete={deleteReservation}
          />
        ))}
        <CreateReservation />
      </div> */}
    </div>
  );
}
