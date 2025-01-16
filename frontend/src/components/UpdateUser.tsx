import axios from 'axios';
import { useState } from 'react';

const UpdateUser: React.FC = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const [updatedUser, setUpdateUser] = useState({
    ID: '',
    FirstName: '',
    LastName: '',
    Email: '',
    Password: '',
  });

  const updateUser = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      await axios.put(`${apiUrl}/api/users/${updatedUser.ID}`, updatedUser);
      setUpdateUser({
        ID: '',
        FirstName: '',
        LastName: '',
        Email: '',
        Password: '',
      });
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        console.error('API Error:', error.response.data.error);
      } else {
        console.error('Error updating user:', error);
      }
    }
  };

  return (
    <div>
      <form
        onSubmit={updateUser}
        className='mb-6 p-4 bg-blue-100 rounded shadow'
      >
        <input
          type='text'
          placeholder='ID'
          value={updatedUser.ID}
          onChange={(e) =>
            setUpdateUser({ ...updatedUser, ID: e.target.value })
          }
          className='mb-2 w-full p-2 border border-gray-300 rounded'
        />
        <input
          type='text'
          placeholder='First Name'
          value={updatedUser.FirstName}
          onChange={(e) =>
            setUpdateUser({ ...updatedUser, FirstName: e.target.value })
          }
          className='mb-2 w-full p-2 border border-gray-300 rounded'
        />
        <input
          type='text'
          placeholder='Last Name'
          value={updatedUser.LastName}
          onChange={(e) =>
            setUpdateUser({ ...updatedUser, LastName: e.target.value })
          }
          className='mb-2 w-full p-2 border border-gray-300 rounded'
        />
        <input
          type='text'
          placeholder='Email'
          value={updatedUser.Email}
          onChange={(e) =>
            setUpdateUser({ ...updatedUser, Email: e.target.value })
          }
          className='mb-2 w-full p-2 border border-gray-300 rounded'
        />
        <button
          className='bg-blue-400 w-full p-2 rounded hover:bg-blue-500'
          type='submit'
        >
          Update User
        </button>
      </form>
    </div>
  );
};

export default UpdateUser;
