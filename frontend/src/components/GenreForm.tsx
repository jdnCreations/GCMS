import axios, { isAxiosError } from 'axios';
import React, { useState } from 'react';

export default function GenreForm() {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const [genre, setGenre] = useState<string>('');

  const handleAddGenre = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (genre == '') return;

    try {
      const response = await axios.post(`${apiUrl}/api/genres`, {
        Name: genre,
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

  return (
    <form
      onSubmit={handleAddGenre}
      className='w-full mb-6 p-4 bg-[#B7A99A] rounded shadow text-gray-800'
    >
      <input
        type='text'
        name='Name'
        onChange={(e) => setGenre(e.target.value)}
        value={genre}
        className='mb-2 w-full p-2 border text-[#4a4a4a] border-gray-300 rounded focus:outline-nook-olive'
        placeholder='Genre'
      />
      <button
        className='bg-nook-olive text-white w-full p-2 rounded hover:bg-nook-light-olive'
        type='submit'
      >
        Add Genre
      </button>
    </form>
  );
}
