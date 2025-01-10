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

const calculateReservationLength = (start: number, end: number) => {
  const startMilli = start / 1000;
  const endMilli = end / 1000;

  const duration = endMilli - startMilli;

  const hours = Math.floor(duration / 3600000);
  const minutes = Math.floor((duration % 3600000) / 60000);

  return {
    hours,
    minutes,
  };
};

const ReservationComponent: React.FC<{
  reservation: Reservation;
  onDelete: (id: string) => void;
}> = ({ reservation, onDelete }) => {
  const handleDelete = () => {
    if (reservation.ID) {
      onDelete(reservation.ID);
    } else {
      console.error('Reservation ID is undefined');
    }
  };

  const reservationLength = calculateReservationLength(
    reservation.StartTime.Microseconds,
    reservation.EndTime.Microseconds
  );

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

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);

    return date.toLocaleDateString('en-AU', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  };

  return (
    <div className='max-w-sm mx-auto bg-white shadow-lg rounded-lg overflow-hidden border border-gray-200 text-black p-2'>
      <h1 className='text-gray-800 font-bold'>{reservation?.GameName}</h1>
      <p className='text-gray-600'>
        Your reservation is for {reservationLength.hours} hours, and{' '}
        {reservationLength.minutes} minutes.
      </p>
      <p>{formatDate(reservation.ResDate)}</p>
      <p>{formatTime(reservation.StartTime.Microseconds)}</p>
      <p>{formatTime(reservation.EndTime.Microseconds)}</p>
      <button
        onClick={handleDelete}
        className='bg-red-800 text-white rounded p-2 font-bold'
      >
        Delete Reservation
      </button>
    </div>
  );
};

export default ReservationComponent;
