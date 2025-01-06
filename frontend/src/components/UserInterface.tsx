'use client';

import React, { useEffect, useState } from 'react';
import CardComponent from './CardComponent';
import axios from 'axios';

interface User {
  ID: string;
  FirstName: string;
  LastName: string;
  Email: string;
}

interface UserInterfaceProps {
  backendName: string;
}

const UserInterface: React.FC<UserInterfaceProps> = ({ backendName }) => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const [users, setUsers] = useState<User[]>([]);
  const [newUser, setNewUser] = useState({
    first_name: '',
    last_name: '',
    email: '',
  });
  const [updatedUser, setUpdateUser] = useState({
    id: '',
    first_name: '',
    last_name: '',
    email: '',
  });

  // fetch all users
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`${apiUrl}/api/users`);
        setUsers(response.data.reverse());
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, [backendName, apiUrl]);

  const createUser = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const response = await axios.post(`${apiUrl}/api/users`, newUser);
      setUsers([response.data, ...users]);
      setNewUser({ first_name: '', last_name: '', email: '' });
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        console.error('API Error:', error.response.data.error);
      } else {
        console.error('Error creating user:', error);
      }
    }
  };

  const updateUser = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const response = await axios.put(
        `${apiUrl}/api/users/${updatedUser.id}`,
        updatedUser
      );
      const updatedUserData = response.data;
      setUsers(
        users.map((user) =>
          user.ID === updatedUserData.ID ? updatedUserData : user
        )
      );
      setUpdateUser({ id: '', first_name: '', last_name: '', email: '' });
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        console.error('API Error:', error.response.data.error);
      } else {
        console.error('Error updating user:', error);
      }
    }
  };

  const deleteUser = async (id: string) => {
    try {
      await axios.delete(`${apiUrl}/api/users/${id}`);
      setUsers(users.filter((user) => user.ID !== id));
    } catch (error) {
      console.error('Error deleting user:', error);
    }
  };

  return (
    <div className='bg-blue-400 justify-center items-center flex flex-col rounded-lg py-2 px-8'>
      {/* Create user */}
      <div>
        <form
          onSubmit={createUser}
          className='mb-6 p-4 bg-blue-100 rounded shadow'
        >
          <input
            type='text'
            placeholder='First Name'
            value={newUser.first_name}
            onChange={(e) =>
              setNewUser({ ...newUser, first_name: e.target.value })
            }
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='Last Name'
            value={newUser.last_name}
            onChange={(e) =>
              setNewUser({ ...newUser, last_name: e.target.value })
            }
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='Email'
            value={newUser.email}
            onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <button
            className='bg-cyan-400 w-full p-2 rounded hover:bg-cyan-500'
            type='submit'
          >
            Create User
          </button>
        </form>
      </div>

      {/* Update user */}
      <div>
        <form
          onSubmit={updateUser}
          className='mb-6 p-4 bg-blue-100 rounded shadow'
        >
          <input
            type='text'
            placeholder='ID'
            value={updatedUser.id}
            onChange={(e) =>
              setUpdateUser({ ...updatedUser, id: e.target.value })
            }
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='First Name'
            value={updatedUser.first_name}
            onChange={(e) =>
              setUpdateUser({ ...updatedUser, first_name: e.target.value })
            }
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='Last Name'
            value={updatedUser.last_name}
            onChange={(e) =>
              setUpdateUser({ ...updatedUser, last_name: e.target.value })
            }
            className='mb-2 w-full p-2 border border-gray-300 rounded'
          />
          <input
            type='text'
            placeholder='Email'
            value={updatedUser.email}
            onChange={(e) =>
              setUpdateUser({ ...updatedUser, email: e.target.value })
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

      <div className='flex flex-col gap-4 justify-center items-center space-y-4'>
        <h1 className='font-bold text-3xl py-4'>Users</h1>
        {users.map((user) => (
          <div key={user.ID} className='flex gap-2 w-full justify-between'>
            <CardComponent card={user} />
            <button
              onClick={() => deleteUser(user.ID)}
              className='bg-orange-400'
            >
              Delete User
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default UserInterface;
