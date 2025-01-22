import { useAuth } from '@/context/AuthContext';
import axios, { isAxiosError } from 'axios';
import React, { useState } from 'react';

interface GameFormData {
  Title: string;
  Copies: number;
}
const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function GameForm() {
  const { jwt } = useAuth();
  const [newGame, setNewGame] = useState<GameFormData>({
    Title: '',
    Copies: 0,
  });

  const handleAddGame = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (
      newGame.Title == '' ||
      !newGame.Copies ||
      isNaN(Number(newGame.Copies))
    ) {
      return;
    }

    try {
      const response = await axios.post(`${apiUrl}/api/games`, newGame, {
        headers: { Authorization: `Bearer ${jwt}` },
      });
      console.log(response.data);
    } catch (error) {
      if (isAxiosError(error) && error) {
        console.error(error.response?.data.error);
      } else {
        console.error(error);
      }
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setNewGame((prev) => ({
      ...prev,
      [name]: name === 'Copies' ? Number(value) : value,
    }));
  };

  return (
    <form
      onSubmit={handleAddGame}
      className='w-full mb-1 p-4 bg-[#B7A99A] rounded shadow text-gray-800'
    >
      <input
        type='text'
        name='Title'
        onChange={handleChange}
        value={newGame.Title}
        className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded focus:outline-nook-olive'
        placeholder='Game Title'
      />
      <input
        type='number'
        name='Copies'
        onChange={handleChange}
        value={newGame.Copies}
        className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded focus:outline-nook-olive'
        placeholder='Copies'
      />
      <button
        className='bg-nook-olive text-white w-full p-2 rounded hover:bg-nook-light-olive'
        type='submit'
      >
        Add Game
      </button>
    </form>
  );
}
