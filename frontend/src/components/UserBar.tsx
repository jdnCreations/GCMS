import { useAuth } from '@/context/AuthContext';
import React from 'react';

export default function UserBar() {
  const { name, logout } = useAuth();
  return (
    <div className='flex justify-between bg-nook-charcoal rounded p-2'>
      <p>Hey, {name}.</p>
      <button className='hover:text-nook-rose' onClick={logout}>
        Log out
      </button>
    </div>
  );
}
