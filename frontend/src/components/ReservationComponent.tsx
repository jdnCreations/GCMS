interface Reservation {
  ID?: string;
  StartTime: string;
  EndTime: string;
  UserID: string;
  GameID: string;
  GameName?: string;
}

const calculateReservationLength = (start: string, end: string) => {
  const startDate = new Date(start).getTime();
  const endDate = new Date(end).getTime();
  const diffTime = endDate - startDate;

  const totalDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

  // Breakdown into days, hours, and minutes
  const days = totalDays;
  const hours = Math.floor(
    (diffTime % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)
  );
  const minutes = Math.floor((diffTime % (1000 * 60 * 60)) / (1000 * 60));

  return { days, hours, minutes };
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
    reservation.StartTime,
    reservation.EndTime
  );

  return (
    <div className='max-w-sm mx-auto bg-white shadow-lg rounded-lg overflow-hidden border border-gray-200 text-black p-2'>
      <h1 className='text-gray-800 font-bold'>{reservation?.GameName}</h1>
      <p className='text-gray-600'>
        Your reservation is for {reservationLength.days} days,{' '}
        {reservationLength.hours} hours, and {reservationLength.minutes}{' '}
        minutes.
      </p>
      <p>
        {new Date(reservation.StartTime).toLocaleString('en-AU', {
          dateStyle: 'medium',
          timeStyle: 'short',
        })}
      </p>
      <p>
        {new Date(reservation.EndTime).toLocaleString('en-AU', {
          dateStyle: 'medium',
          timeStyle: 'short',
        })}
      </p>
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
