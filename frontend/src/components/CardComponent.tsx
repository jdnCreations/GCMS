'use client';
import React from 'react';

interface Card {
  ID: string;
  FirstName: string;
  LastName: string;
  Email: string;
}

const CardComponent: React.FC<{ card: Card }> = ({ card }) => {
  return (
    <div className='flex flex-col gap-2 p-8 bg-cyan-400 rounded-sm'>
      <p className='font-thin'>{card.ID}</p>
      <h2 className='font-bold'>{card.FirstName}</h2>
      <p>{card.Email}</p>
    </div>
  );
};

export default CardComponent;
