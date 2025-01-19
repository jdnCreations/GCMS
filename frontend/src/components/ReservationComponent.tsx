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
    <div className='max-w-sm w-full mx-auto bg-white shadow-lg rounded-lg overflow-hidden border border-gray-200 text-black p-2'>
      <div className='text-nook-charcoal'>
        You have reserved the game{' '}
        <span className='font-bold'>{reservation?.GameName}</span> for{' '}
        {reservationLength.hours} hours, and {reservationLength.minutes}{' '}
        minutes.
      </div>
      <p>{formatDate(reservation.ResDate)}</p>
      <div className='flex justify-between'>
        <p>{formatTime(reservation.StartTime.Microseconds)}</p>
        <p>to</p>
        <p>{formatTime(reservation.EndTime.Microseconds)}</p>
      </div>
      <button
        onClick={handleDelete}
        className='bg-nook-dark-rose text-white rounded p-2 my-1 w-full hover:bg-nook-rose'
      >
        Delete Reservation
      </button>
    </div>
  );
};

export default ReservationComponent;
